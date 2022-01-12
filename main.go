package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"unsafe"
)

type VM struct {
	stack   []int64
	maxSize uint
	bp      uint // stack base pointer
	sp      uint // stack pointer
	ip      uint // instruction pointer
	stop    bool
	program *Program
}

const PROGRAM_CAPACITY = 512

type Program = []Inst

func NewVM(size uint, program *Program) *VM {
	return &VM{
		stack:   make([]int64, size, size),
		maxSize: size,
		program: program,
	}
}

type InstKind int64

type Inst struct {
	Kind    InstKind
	Operand int64
}

const HALT_VALUE = 0x12345

const (
	Inst_Push InstKind = iota
	Inst_Add
	Inst_Sub
	Inst_Mul
	Inst_Div
	Inst_Eq
	Inst_Halt
	Inst_Jmp
	Inst_JmpTrue
	Inst_Dup
	Inst_Print
	Inst_Count
)

func NewInst(kind InstKind) func(int64) Inst {
	return func(value int64) Inst {
		return Inst{Kind: kind, Operand: value}
	}
}

var (
	// no operand
	Add   = Inst{Kind: Inst_Add}   // add
	Sub   = Inst{Kind: Inst_Sub}   // substract
	Mul   = Inst{Kind: Inst_Mul}   // multiply
	Div   = Inst{Kind: Inst_Div}   // divide
	Print = Inst{Kind: Inst_Print} // print the value at the top of the stack
	Halt  = Inst{Kind: Inst_Halt}  // stop the vm
	Eq    = Inst{Kind: Inst_Eq}    // check if last 2 values are equal and but 1 or 0 at the top

	//operand
	Push    = NewInst(Inst_Push)    // push operand at the top of the stack
	Jmp     = NewInst(Inst_Jmp)     // Jmp at a position of the program
	JmpTrue = NewInst(Inst_JmpTrue) // Jump if top value of the stack != 0 at the position of the program
	Dup     = NewInst(Inst_Dup)     // Duplicate the value at the relative position in stack at the top of the stack
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
	case Inst_Eq:
		return "EQ"
	case Inst_Halt:
		return "HALT"
	case Inst_Jmp:
		return "JMP"
	case Inst_JmpTrue:
		return "JMP_TRUE"
	case Inst_Dup:
		return "DUP"
	case Inst_Print:
		return "PRINT"
	default:
		panic(fmt.Errorf("InstKind unkown human representation of error: %d", ik))
	}

}

func (i Inst) String() string {
	switch i.Kind {
	case Inst_Push, Inst_Jmp, Inst_JmpTrue, Inst_Dup:
		return fmt.Sprintf("%v[%v]", i.Kind, i.Operand)
	case Inst_Add, Inst_Halt, Inst_Sub, Inst_Mul, Inst_Div, Inst_Print:
		return fmt.Sprintf("%v", i.Kind)
	default:
		panic(fmt.Errorf("Inst unkown human representation of error: %d", i.Kind))
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

func (v *VM) executeInst(inst Inst) (err Err) {
	switch inst.Kind {
	case Inst_Push:
		if v.sp >= v.maxSize {
			err = Err_Overflow
		} else {
			v.stack[v.sp] = inst.Operand
			v.sp++
			v.ip++
		}
	case Inst_Add, Inst_Sub, Inst_Mul, Inst_Div, Inst_Eq:
		if len(v.stack) < 2 {
			err = Err_Underflow
		} else {
			res, err2 := binaryOp(v.stack[v.sp-2], v.stack[v.sp-1], inst.Kind)
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
		if inst.Operand < 0 || inst.Operand >= int64(len(*v.program)) {
			err = Err_OutOfIndexInstruction
		} else {

			v.ip = uint(inst.Operand)
		}
	case Inst_JmpTrue:
		if inst.Operand < 0 || inst.Operand >= int64(len(*v.program)) {
			err = Err_OutOfIndexInstruction
		} else if v.stack[v.sp-1] != 0 {
			v.ip = uint(inst.Operand)
		} else {
			v.ip++
		}
	case Inst_Dup: // duplicate relative to sp
		if inst.Operand <= 0 {
			err = Err_OutOfIndexInstruction
		} else {
			v.stack[v.sp] = v.stack[v.sp-uint(inst.Operand)]
			v.sp++
			v.ip++
		}
	case Inst_Print:
		fmt.Println(v.stack[v.sp-1])
		v.ip++
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

// func NewProgram(insts ...Inst) *Program {
// 	i := 0
// 	ret := new(Program)
// 	for _, inst := range insts {
// 		(*ret)[i] = inst
// 		i++
// 	}
// 	fmt.Println(*ret)
// 	return ret
// }

func NewProgram(insts ...Inst) *Program {
	ret := make([]Inst, 0, len(insts))
	for _, inst := range insts {
		ret = append(ret, inst)
	}
	return &ret
}
func (v *VM) execute() {
	counter := 0
	for !v.stop && counter < 100 {
		inst := (*v.program)[v.ip]
		// fmt.Printf("inst=%v,ip=%v, sp=%v\n", inst, v.ip, v.sp)
		err := v.executeInst(inst)
		if err != OK {
			panic(err)
		}
		// v.dump()
		counter++
	}
}

func (v *VM) WriteToFile(pathFile string) error {
	fd, err := os.Create(pathFile)
	if err != nil {
		return err
	}
	defer fd.Close()
	return binary.Write(fd, binary.BigEndian, v.program)
}

func LoadProgram(pathFile string) (*Program, error) {
	fd, err := os.Open(pathFile)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	sizeInst := int64(unsafe.Sizeof(Inst{}))

	fi, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	sizeFile := fi.Size()

	fmt.Println(unsafe.Sizeof(uintptr(0)))
	fmt.Println(unsafe.Sizeof(int64(0)))

	p := make([]Inst, sizeFile/sizeInst, sizeFile/sizeInst)
	err = binary.Read(fd, binary.BigEndian, &p)
	if err != nil {
		return nil, err
	}
	fmt.Println(p)
	return &p, nil
}

func main() {

	code, err := ioutil.ReadFile("./toto.evm")
	if err != nil {
		Panic("could not read file: %v", err)
	}

	p := loadSourceCode(string(code))
	// p, err := LoadProgram("./toto.vm")
	// if err != nil {
	// 	panic(err)
	// }
	vm := NewVM(512, p)
	vm.execute()
	// err := vm.WriteToFile("./toto.vm")
	// if err != nil {
	// 	panic(err)
	// }
}
