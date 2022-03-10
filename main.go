package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fmarmol/basename/pkg/basename"
	"gopkg.in/alecthomas/kingpin.v2"
)

func Panic(s string, args ...interface{}) {
	err := fmt.Errorf(s, args...)
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}

const MEM_CAPACITY = 1024

type VM struct {
	mem     [MEM_CAPACITY]byte
	stack   []Word
	maxSize uint32
	bp      uint32 // stack base pointer
	sp      uint32 // stack pointer
	ip      uint32 // instruction pointer
	stop    bool
	program *Program
}

func (v *VM) stackTop() Word {
	return v.stack[v.sp-1]
}

const PROGRAM_CAPACITY = 512
const STACK_CAPACITY = 100

type Program []Inst

func (p *Program) Size() uint32 {
	return uint32(len(*p))
}

// size is the max size of the stack
func NewVM(size uint32, program *Program) *VM {
	return &VM{
		stack:   make([]Word, size, size),
		maxSize: size,
		program: program,
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
	Err_WrongTypeOperation
	Err_SpaceNotFound
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
	case Err_WrongTypeOperation:
		return "Wrong Type Operation"
	case Err_SpaceNotFound:
		return "Not enough space to allocate memory"
	default:
		Panic("Err unknown human representation of error: %d", e)
	}
	return ""
}

func NewArray(size int, value byte) []byte {

	ret := make([]byte, 0, size)

	for i := 0; i < size; i++ {
		ret = append(ret, value)
	}
	return ret
}

func (v *VM) freeMemSpace(ptr uint32, size uint32) {
	copy(v.mem[int(ptr):int(ptr+size)], NewArray(int(size), 0))
}

func (v *VM) findMemSpace(size uint32) (uint32, Err) { // return the index on the memory able to allocate size bytes
LOOP:
	for i := range v.mem {
		for j := range v.mem[i : i+int(size)] {
			if v.mem[i+j] != 0 {
				continue LOOP
			}
		}
		return uint32(i), OK
	}
	return 0, Err_SpaceNotFound
}

func (v *VM) executeInst(inst Inst) (err Err) {
	switch inst.Kind {
	case Inst_Alloc:
		nbytes := v.stack[v.sp-1].UInt32()
		v.stack[v.sp-1] = Word{}
		v.sp--
		addr, err2 := v.findMemSpace(nbytes)
		if err2 != OK {
			err = err2
		} else {
			v.stack[v.sp] = NewWord(addr, Ptr)
			v.sp++
		}
		v.ip++
	case Inst_Start:
		v.ip++
	case Inst_PushInt, Inst_PushFloat, Inst_PushUInt32:
		if v.sp >= v.maxSize {
			err = Err_Overflow
		} else {
			v.stack[v.sp] = inst.Operand
			v.sp++
			v.ip++
		}
	// DEBUG INSTRUCTIONS
	case Inst_Dump:
		v.dump()
		v.ip++
	case Inst_EqInt:
		a := v.stack[v.sp-1]
		op := inst.Operand
		if a.Kind != op.Kind {
			fmt.Fprintf(os.Stderr, "incompatible types: tried comparison between %v and %v\n", a.Kind, op.Kind)
			err = Err_WrongTypeOperation
			return
		}
		if a.Int64() != op.Int64() {
			Panic("invalid assertion top[%v] != eq[%v]", a.Int64(), op.Int64())
		}
		v.ip++
	case Inst_EqFloat:
		top := v.stackTop()
		op := inst.Operand
		if top.Kind != op.Kind {
			fmt.Fprintf(os.Stderr, "incompatible types: tried comparison between %v and %v\n", top.Kind, op.Kind)
			err = Err_WrongTypeOperation
			return
		}
		if top.Float64() != op.Float64() {
			Panic("invalid assertion top[%v] != eq[%v]", top.Float64(), op.Float64())
		}
		v.ip++
	// END DEBUG INSTRUCTIONS
	case Inst_Add, Inst_Sub, Inst_Mul, Inst_Div, Inst_Eq:
		if len(v.stack) < 2 {
			err = Err_Underflow
			return
		}
		a := v.stack[v.sp-2]
		b := v.stack[v.sp-1]
		if a.Kind != b.Kind {
			fmt.Fprintf(os.Stderr, "incompatible types: tried to binary operation between %v and %v\n", a.Kind, b.Kind)
			err = Err_WrongTypeOperation
			return
		}
		var result Word
		switch a.Kind {
		case Int64:
			res, err2 := binaryOp(a.Int64(), b.Int64(), inst.Kind)
			if err2 != OK {
				err = err2
				return
			}
			result = NewWord(res, Int64)
		case Float64:
			res, err2 := binaryOp(a.Float64(), b.Float64(), inst.Kind)
			if err2 != OK {
				err = err2
				return
			}
			result = NewWord(res, Float64)
		}
		v.stack[v.sp-2] = result
		v.stack[v.sp-1] = Word{}
		v.sp--
		v.ip++
	case Inst_Swap:
		pos_top := v.sp - 1
		pos_sec := v.sp - (inst.Operand.UInt32())
		v.stack[pos_sec], v.stack[pos_top] = v.stack[pos_top], v.stack[pos_sec]
		v.ip++
	case Inst_Drop:
		v.stack[v.sp-1] = Word{}
		v.sp--
		v.ip++
	case Inst_Halt:
		v.stop = true
	case Inst_Ret:
		v.ip = v.stackTop().UInt32()
		v.stack[v.sp-1] = Word{}
		v.sp--
	case Inst_Call:
		if inst.Operand.UInt32() < 0 || inst.Operand.UInt32() >= v.program.Size() {
			err = Err_OutOfIndexInstruction
		} else {
			v.stack[v.sp] = NewWord(v.ip+1, UInt32)
			v.sp++
			v.ip = inst.Operand.UInt32()
		}

	case Inst_Jmp:
		if inst.Operand.UInt32() < 0 || inst.Operand.UInt32() >= v.program.Size() {
			err = Err_OutOfIndexInstruction
		} else {
			v.ip = inst.Operand.UInt32()
		}
	case Inst_JmpTrue:
		if inst.Operand.UInt32() < 0 || inst.Operand.UInt32() >= v.program.Size() {
			err = Err_OutOfIndexInstruction
		} else if !v.stack[v.sp-1].IsZero() {
			v.ip = inst.Operand.UInt32()
		} else {
			v.ip++
		}
	case Inst_Dup: // duplicate relative to sp
		if inst.Operand.UInt32() <= 0 {
			err = Err_OutOfIndexInstruction
		} else {
			v.stack[v.sp] = v.stack[v.sp-inst.Operand.UInt32()]
			v.sp++
			v.ip++
		}
	case Inst_Print:
		fmt.Println("->", v.stack[v.sp-1])
		v.ip++
	case Inst_Label:
		v.ip++
	default:
		err = Err_IllegalInstruction
	}
	return
}

func (v *VM) dump() {
	fmt.Println("STACK:")
	for i := v.bp; i < v.sp; i++ {
		word := v.stack[i]
		switch word.Kind {
		case Int64:
			fmt.Printf("\t addr=%v %v %v\n", i, word.Kind, word.Int64())
		case UInt32:
			fmt.Printf("\t addr=%v %v %v\n", i, word.Kind, word.UInt32())
		case Float64:
			fmt.Printf("\t addr=%v %v %v\n", i, word.Kind, word.Float64())
		case Ptr:
			fmt.Printf("\t addr=%v %v %v\n", i, word.Kind, word.Ptr())
		default:
			Panic("cannot dump word of kind: %v, %v", word.Kind, v.sp)

		}
	}
}

func NewProgram(insts ...Inst) *Program {
	ret := make([]Inst, 0, len(insts))
	for _, inst := range insts {
		ret = append(ret, inst)
	}
	p := Program(ret)
	return &p
}

func (v *VM) execute(maxStep uint) {
	var counter uint
	var started bool
	for !v.stop && counter < maxStep {
		inst := (*v.program)[v.ip]
		if inst.Kind != Inst_Start && !started {
			v.ip++
			continue
		} else {
			started = true
		}
		fmt.Printf("inst=%v,ip=%v, sp=%v\n", inst, v.ip, v.sp)
		err := v.executeInst(inst)
		if err != OK {
			fmt.Fprintf(os.Stderr, "Inst: %v, Err: %v\n", inst.String(), err.String())
			v.dump()
			return
		}
		counter++
	}
	fmt.Println("number of execution steps:", counter)
}

func (v *VM) WriteToFile(pathFile string) error {
	fd, err := os.Create(pathFile)
	if err != nil {
		return err
	}
	defer fd.Close()
	return binary.Write(fd, binary.BigEndian, *v.program)
}

var (
	app    = kingpin.New("vm", "vm main command")
	comp   = app.Command("compile", "compile a .evm file").Alias("c")
	source = comp.Arg("source", "source file").String()
	output = comp.Flag("output", "output file .vm").Short('o').String()

	run       = app.Command("run", "run vm file").Alias("r")
	sourceRun = run.Arg("source", "source file .vm").String()
	maxStep   = run.Flag("max_step", "max exection steps allowed").Default("300").Uint()

	disas       = app.Command("disas", "disassemble a program .vm")
	sourceDisas = disas.Arg("source", "source file .vm").String()
	outputDisas = disas.Flag("output", "output file .vm.disas").Short('o').String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case comp.FullCommand():
		fi := basename.ParseFile(*source)
		code, err := ioutil.ReadFile(*source)
		if err != nil {
			Panic("could not read file: %v", err)
		}
		p := loadSourceCode(string(code))
		vm := NewVM(PROGRAM_CAPACITY, p)
		path := filepath.Join(fi.Dir, fi.Basename) + ".vm"
		err = vm.WriteToFile(path)
		if err != nil {
			panic(err)
		}
	case run.FullCommand():
		p, err := LoadProgram(*sourceRun)
		if err != nil {
			panic(err)
		}
		vm := NewVM(PROGRAM_CAPACITY, p)
		vm.execute(*maxStep)
	case disas.FullCommand():
		p, err := LoadProgram(*sourceDisas)
		if err != nil {
			panic(err)
		}
		ps := p.disas()
		if outputDisas == nil {
			for _, inst := range ps {
				fmt.Println(inst)
			}
		} else {
			fd, err := os.Create(*outputDisas)
			if err != nil {
				panic(err)
			}
			defer fd.Close()
			content := strings.Join(ps, "\n")
			fd.WriteString(content)
		}
	}
}
