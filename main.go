package main

import "fmt"

type VM struct {
	stack   []int
	maxSize uint
	bp      uint // stack base pointer
	sp      uint // stack pointer
	ip      uint // instruction pointer
	stop    bool
	program *Program
}

type Program []Inst

func NewVM(size uint, program *Program) *VM {
	return &VM{
		stack:   make([]int, size, size),
		maxSize: size,
		program: program,
	}
}

type InstKind int

type Inst struct {
	kind    InstKind
	operand int
}

const HALT_VALUE = 0x12345

const (
	Inst_Push InstKind = iota
	Inst_Add
	Inst_Sub
	Inst_Mul
	Inst_Div
	Inst_Halt
	Inst_Jmp
	Inst_JmpTrue
)

var (
	Add  = Inst{kind: Inst_Add}
	Sub  = Inst{kind: Inst_Sub}
	Mul  = Inst{kind: Inst_Mul}
	Div  = Inst{kind: Inst_Div}
	Push = func(value int) Inst {
		return Inst{kind: Inst_Push, operand: value}
	}
	Halt = Inst{kind: Inst_Halt}
	Jmp  = func(value int) Inst {
		return Inst{kind: Inst_Jmp, operand: value}
	}
	JmpTrue = func(value int) Inst {
		return Inst{kind: Inst_JmpTrue, operand: value}
	}
)

func (ik InstKind) String() string {
	switch ik {
	case Inst_Push:
		return "PUSH"
	case Inst_Add:
		return "ADD"
	case Inst_Sub:
		return "SUB"
	case Inst_Mul:
		return "MUL"
	case Inst_Div:
		return "DIV"
	case Inst_Halt:
		return "HALT"
	case Inst_Jmp:
		return "JMP"
	case Inst_JmpTrue:
		return "JMP_TRUE"
	default:
		panic(fmt.Errorf("InstKind unkown human representation of error: %d", ik))
	}

}

func (i Inst) String() string {
	switch i.kind {
	case Inst_Push, Inst_Jmp, Inst_JmpTrue:
		return fmt.Sprintf("%v[%v]", i.kind, i.operand)
	case Inst_Add, Inst_Halt, Inst_Sub, Inst_Mul, Inst_Div:
		return fmt.Sprintf("%v", i.kind)
	default:
		panic(fmt.Errorf("Inst unkown human representation of error: %d", i.kind))
	}
}

type Err int

const (
	OK Err = iota
	Err_Overflow
	Err_Underflow
	Err_IllegalInstruction
	Err_DivisionByZero
	Err_OutOfIndexInstruction
)

func (e Err) String() string {
	switch e {
	case Err_Overflow:
		return "ERROR OVERFLOW"
	case Err_Underflow:
		return "ERORR UNDERFLOW"
	case Err_IllegalInstruction:
		return "ERROR ILLEGAL INSTRUCTION"
	case Err_DivisionByZero:
		return "Division By Zero"
	case OK:
		return "OK"
	case Err_OutOfIndexInstruction:
		return "Out Of Index Instruction"
	default:
		panic(fmt.Errorf("Err unkown human representation of error: %d", e))
	}
}

func binaryOp(arg1 int, arg2 int, ik InstKind) (ret int, err Err) {
	switch ik {
	case Inst_Add:
		ret = arg1 + arg2
	case Inst_Sub:
		ret = arg1 - arg2
	case Inst_Mul:
		ret = arg1 * arg2
	case Inst_Div:
		if arg2 == 0 {
			err = Err_DivisionByZero
		} else {
			ret = arg1 / arg2
		}
	default:
		panic("unknown binaryOp")
	}
	return
}

func (v *VM) executeInst(inst Inst) (err Err) {
	switch inst.kind {
	case Inst_Push:
		if v.sp >= v.maxSize {
			err = Err_Overflow
		} else {
			v.stack[v.sp] = inst.operand
			v.sp++
			v.ip++
		}
	case Inst_Add, Inst_Sub, Inst_Mul, Inst_Div:
		if len(v.stack) < 2 {
			err = Err_Underflow
		} else {
			res, err2 := binaryOp(v.stack[v.sp-2], v.stack[v.sp-1], inst.kind)
			if err2 == OK {
				v.stack[v.sp-2] = res
				v.stack[v.sp-1] = 0
				v.sp--
				v.ip++
			}
		}
	case Inst_Halt:
		v.stop = true
	case Inst_Jmp:
		if inst.operand < 0 || inst.operand >= len(*v.program) {
			err = Err_OutOfIndexInstruction
		} else {

			v.ip = uint(inst.operand)
		}
	case Inst_JmpTrue:
		if inst.operand < 0 || inst.operand >= len(*v.program) {
			err = Err_OutOfIndexInstruction
		} else if v.stack[v.sp-1] != 0 {
			v.ip = uint(inst.operand)
		} else {
			v.ip++
		}
	default:
		err = Err_IllegalInstruction
	}
	return
}

func (v *VM) dump() {
	fmt.Println("STACK:")
	for i := v.bp; i < v.sp; i++ {
		fmt.Println("\t", v.stack[i])
	}
}

func NewProgram(size uint, insts ...Inst) *Program {
	ret := new(Program)
	*ret = append(*ret, insts...)
	return ret
}

func (v *VM) execute() {
	counter := 0
	for !v.stop && counter < 50 {
		inst := (*v.program)[v.ip]
		fmt.Printf("inst=%v,ip=%v\n", inst, v.ip)
		err := v.executeInst(inst)
		if err != OK {
			panic(err)
		}
		v.dump()
		counter++
	}
}

func main() {
	p := NewProgram(10,
		Push(10), // 0
		Push(1),  // 1
		Sub,      // 2
		JmpTrue(1),
		Halt,
	)
	vm := NewVM(10, p)
	vm.execute()
}
