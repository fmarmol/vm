package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/fmarmol/basename/pkg/basename"
	"gopkg.in/alecthomas/kingpin.v2"
)

func Panic(s string, args ...interface{}) {
	err := fmt.Errorf(s, args...)
	panic(err)
}

type VM struct {
	stack   []Word
	maxSize uint32
	bp      uint32 // stack base pointer
	sp      uint32 // stack pointer
	ip      uint32 // instruction pointer
	stop    bool
	program *Program
}

const PROGRAM_CAPACITY = 512
const STACK_CAPACITY = 100

type Program []Inst

func (p *Program) Size() uint32 {
	return uint32(len(*p))
}

func NewVM(size uint32, program *Program) *VM {
	return &VM{
		stack:   make([]Word, size, size),
		maxSize: size,
		program: program,
	}
}

type InstKind uint32

type Inst struct {
	Kind    InstKind
	Operand Word
}

const (
	Inst_PushInt InstKind = iota
	Inst_PushFloat
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
	Inst_Label
	Inst_Com
	Inst_Count
)

func NewInst(kind InstKind) func(Word) Inst {
	return func(value Word) Inst {
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
	PushInt   = NewInst(Inst_PushInt)   // push integer at the top of the stack
	PushFloat = NewInst(Inst_PushFloat) // push float at the top of the stack
	Jmp       = NewInst(Inst_Jmp)       // Jmp at a position of the program
	JmpTrue   = NewInst(Inst_JmpTrue)   // Jump if top value of the stack != 0 at the position of the program
	Dup       = NewInst(Inst_Dup)       // Duplicate the value at the relative position in stack at the top of the stack
	Label     = NewInst(Inst_Label)
)

func (ik InstKind) String() string {
	switch ik {
	case Inst_PushInt:
		return "pushi"
	case Inst_PushFloat:
		return "pushf"
	case Inst_Add:
		return "add"
	case Inst_Sub:
		return "sub"
	case Inst_Mul:
		return "mul"
	case Inst_Div:
		return "div"
	case Inst_Eq:
		return "eq"
	case Inst_Halt:
		return "halt"
	case Inst_Jmp:
		return "jmp"
	case Inst_JmpTrue:
		return "jmptrue"
	case Inst_Dup:
		return "dup"
	case Inst_Print:
		return "print"
	case Inst_Label:
		return "label"
	case Inst_Com:
		return "comment"
	default:
		Panic("InstKind unknown human representation of error: %d", ik)
	}
	return ""
}

func (i Inst) String() string {
	switch i.Kind {
	case Inst_PushInt, Inst_PushFloat, Inst_Jmp, Inst_JmpTrue, Inst_Dup, Inst_Label:
		return fmt.Sprintf("%v %v", i.Kind, i.Operand)
	case Inst_Add, Inst_Halt, Inst_Sub, Inst_Mul, Inst_Div, Inst_Print:
		return fmt.Sprintf("%v", i.Kind)
	default:
		Panic("Inst unknown human representation of error: %d", i.Kind)
	}
	return ""
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
		Panic("Err unkown human representation of error: %d", e)
	}
	return ""
}

func (v *VM) executeInst(inst Inst) (err Err) {
	switch inst.Kind {
	case Inst_PushInt:
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
			return
		}
		a := v.stack[v.sp-2]
		b := v.stack[v.sp-1]
		if a.Kind != b.Kind {
			panic(fmt.Errorf("incompatible types: tried to to binary operation between %v and %v", a.Kind, b.Kind))
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

	case Inst_Halt:
		v.stop = true
	case Inst_Jmp:
		if inst.Operand.UInt32() < 0 || inst.Operand.UInt32() >= v.program.Size() {
			err = Err_OutOfIndexInstruction
		} else {
			v.ip = inst.Operand.UInt32()
		}
	case Inst_JmpTrue:
		if inst.Operand.UInt32() < 0 || inst.Operand.UInt32() >= v.program.Size() {
			err = Err_OutOfIndexInstruction
		} else if v.stack[v.sp-1].UInt32() != 0 {
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
		fmt.Println(v.stack[v.sp-1])
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
		fmt.Println("\t", v.stack[i])
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

func (v *VM) execute() {
	counter := 0
	for !v.stop && counter < 100 {
		inst := (*v.program)[v.ip]
		// fmt.Printf("inst=%v,ip=%v, sp=%v\n", inst, v.ip, v.sp)
		err := v.executeInst(inst)
		if err != OK {
			Panic("Inst: %v, Err: %v", inst.String(), err.String())
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

	p := Program(make([]Inst, sizeFile/sizeInst, sizeFile/sizeInst))
	err = binary.Read(fd, binary.BigEndian, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

var (
	app    = kingpin.New("vm", "vm main command")
	comp   = app.Command("compile", "compile a .evm file")
	source = comp.Arg("source", "source file").String()
	output = comp.Flag("output", "output file .vm").Short('o').String()

	run       = app.Command("run", "run vm file")
	sourceRun = run.Arg("source", "source file .vm").String()

	disas       = app.Command("disas", "disassemble a program .vm")
	sourceDisas = disas.Arg("source", "source file .vm").String()
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
		vm.execute()
	case disas.FullCommand():
		// TODO: disasembly output should be able to be input for compile command
		p, err := LoadProgram(*sourceDisas)
		if err != nil {
			panic(err)
		}
		ps := p.disas()
		for _, inst := range ps {
			fmt.Println(inst)
		}
	}
}
