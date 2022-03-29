package vm

import (
	"errors"
	"fmt"

	"github.com/fmarmol/vm/pkg/fatal"
	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/procs"
	"github.com/fmarmol/vm/pkg/rorre"
	"github.com/fmarmol/vm/pkg/word"
)

func (v *VM) IP() uint32          { return v.ip }
func (v *VM) ProgramSize() uint32 { return v.Program.Size() }
func (v *VM) SP() uint32          { return v.sp }
func (v *VM) StackCap() uint32    { return uint32(len(v.Stack)) }
func (v *VM) StackPush(w word.Word) error {
	if v.sp >= v.StackCap() {
		return rorre.Err_Overflow
	}
	v.Stack[v.sp] = w
	v.sp++
	return nil
}

func (v *VM) Stop() {
	v.stop = true
}

func (v *VM) Swap(first, second uint32) error {
	if first < 1 {
		return errors.New("first index must be greater or equal to 1, its relative index. stack[sp-first]")
	}
	if v.sp-first < 0 {
		return errors.New("tried to acces negative index from the stack. stack[sp-first]")
	}
	if second < 1 {
		return errors.New("second index must be greater or equal to 1, its relative index. stack[sp-second]")
	}
	if v.sp-second < 0 {
		return errors.New("tried to acces negative index from the stack. stack[sp-second]")
	}
	v.Stack[v.sp-first], v.Stack[v.sp-second] = v.Stack[v.sp-second], v.Stack[v.sp-first]
	return nil
}

func (v *VM) stackTop() word.Word {
	return v.Stack[v.sp-1]
}

func (v *VM) StackPeek() word.Word {
	if v.sp-1 < 0 {
		panic("need to implement proper error management")
	}
	return v.Stack[v.sp-1]
}

func (v *VM) StackPeekIndex(index uint32) (word.Word, error) {
	if v.sp-index < 0 {
		return word.Word{}, errors.New("need to implement proper error management")
	}
	return v.Stack[v.sp-index], nil
}

func (v *VM) StackPop() (word.Word, error) {
	if v.sp < 1 {
		return word.Word{}, rorre.Err_Underflow
	}
	top := v.StackPeek()
	v.Stack[v.sp-1] = word.Word{}
	v.sp--
	return top, nil

}

// size is the max size of the stack
func NewVM(innerVM InnerVM) *VM {
	v := &VM{}
	v.Program = innerVM.Program
	v.Memory = innerVM.Memory
	return v
}

func (v *VM) Execute(maxStep uint) {
	var counter uint
	var started bool

	rules := loadRulesProcs()
	for !v.stop && counter < maxStep {
		_inst := v.Program[v.ip]
		if _inst.Kind != inst.Inst_Start && !started {
			v.ip++
			continue
		} else {
			started = true
		}
		rule, ok := rules[_inst.Kind]
		if !ok {
			fatal.Panic("rule %v not implemented", _inst.Kind)
		}
		err := rule.proc(v, _inst)
		if err != nil {
			panic(err)
		}
		rule.fip(&IpExec{vm: v, _inst: _inst})
		counter++
	}
	fmt.Println("number of execution steps:", counter)
}

func (v *VM) ExecuteWithDebug(maxStep uint) {
	var counter uint
	var started bool

	rules := loadRulesProcs()
	for !v.stop && counter < maxStep {
		_inst := v.Program[v.ip]
		if _inst.Kind != inst.Inst_Start && !started {
			v.ip++
			continue
		} else {
			started = true
		}
		fmt.Printf("inst=%v,ip=%v, sp=%v\n", _inst, v.ip, v.sp)
		rule, ok := rules[_inst.Kind]
		if !ok {
			fatal.Panic("rule %v not implemented", _inst.Kind)
		}
		err := rule.proc(v, _inst)
		if err != nil {
			panic(err)
		}
		rule.fip(&IpExec{vm: v, _inst: _inst})
		v.dump()
		fmt.Scanln()
		counter++
	}
	fmt.Println("number of execution steps:", counter)
}

func (v *VM) dump() {
	fmt.Println("STACK:")
	for i := v.bp; i < v.sp; i++ {
		_word := v.Stack[i]
		switch _word.Kind {
		case word.Int64:
			fmt.Printf("\t addr=%v %v %v\n", i, _word.Kind, _word.Int64())
		case word.UInt32:
			fmt.Printf("\t addr=%v %v %v\n", i, _word.Kind, _word.UInt32())
		case word.Float64:
			fmt.Printf("\t addr=%v %v %v\n", i, _word.Kind, _word.Float64())
		case word.Ptr:
			fmt.Printf("\t addr=%v %v %v\n", i, _word.Kind, _word.Ptr())
		default:
			fatal.Panic("cannot dump word of kind: %v, %v", _word.Kind, v.sp)

		}
	}
}

// TODO: check ip boundaries
func incIp(ipExec *IpExec) {
	ipExec.vm.ip++
}

func nopIp(*IpExec) {}

func retIp(ipExec *IpExec) {
	top, err := ipExec.vm.StackPop()
	if err != nil {
		panic(err)
	}
	ipExec.vm.ip = top.UInt32()
}

// TODO: check ip boundaries
func jmpIp(ipExec *IpExec) {
	ipExec.vm.ip = ipExec._inst.Operand.UInt32()
}

// TODO: check ip boundaries
func jmpTrueIp(ipExec *IpExec) {
	if !ipExec.vm.StackPeek().IsZero() {
		ipExec.vm.ip = ipExec._inst.Operand.UInt32()
	} else {
		ipExec.vm.ip++
	}

}

// TODO: check ip boundaries
func callIp(ipExec *IpExec) {
	ipExec.vm.ip = ipExec._inst.Operand.UInt32()
}

type IpExec struct {
	vm    *VM
	_inst inst.Inst
}

type ProcExec struct {
	proc func(v procs.VMer, _inst inst.Inst) error
	fip  func(ipExec *IpExec)
}

// func loadRules() map[inst.InstKind](func(v procs.VMer, _inst inst.Inst) error) {
// 	return map[inst.InstKind](func(v procs.VMer, _inst inst.Inst) error){
// 		inst.Inst_Start:      procs.Start,
// 		inst.Inst_PushInt:    procs.Push,
// 		inst.Inst_PushFloat:  procs.Push,
// 		inst.Inst_PushUInt32: procs.Push,
// 	}
// }

// TODO: procs return IpExec and Error to separate properly responsabilities of stack and ip

func loadRulesProcs() map[inst.InstKind]ProcExec {
	return map[inst.InstKind]ProcExec{
		inst.Inst_Start:      {procs.Start, incIp},
		inst.Inst_Label:      {procs.Nop, incIp},
		inst.Inst_PushInt:    {procs.Push, incIp},
		inst.Inst_PushFloat:  {procs.Push, incIp},
		inst.Inst_PushUInt32: {procs.Push, incIp},
		inst.Inst_EqFloat:    {procs.Eq, incIp},
		inst.Inst_EqInt:      {procs.Eq, incIp},
		inst.Inst_Add:        {procs.Bin, incIp},
		inst.Inst_Sub:        {procs.Bin, incIp},
		inst.Inst_Mul:        {procs.Bin, incIp},
		inst.Inst_Div:        {procs.Bin, incIp},
		inst.Inst_Eq:         {procs.Bin, incIp},
		inst.Inst_Swap:       {procs.Swap, incIp},
		inst.Inst_Drop:       {procs.Drop, incIp},
		inst.Inst_Halt:       {procs.Stop, incIp},
		inst.Inst_Ret:        {procs.Nop, retIp}, // special case where need to consume from ip func
		inst.Inst_Call:       {procs.Call, callIp},
		inst.Inst_Jmp:        {procs.Nop, jmpIp},
		inst.Inst_JmpTrue:    {procs.Nop, jmpTrueIp},
		inst.Inst_Dup:        {procs.Dup, incIp},
		inst.Inst_Print:      {procs.Print, incIp},
		inst.Inst_PrintChar:  {procs.PrintChar, incIp},
		inst.Inst_Debug:      {procs.Debug, incIp},
	}
}
