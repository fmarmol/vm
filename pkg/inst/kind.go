package inst

import "github.com/fmarmol/vm/pkg/fatal"

type InstKind uint32

const (
	Inst_PushInt InstKind = iota + 1 // default value 0 should not be valid for easier debug life
	Inst_Push
	Inst_PushFloat
	Inst_PushUInt32
	Inst_EqInt
	Inst_EqFloat
	Inst_Add
	Inst_Sub
	Inst_Mul
	Inst_Div
	Inst_Eq
	Inst_Ret
	Inst_Halt
	Inst_Call
	Inst_Jmp
	Inst_JmpTrue
	Inst_Dup
	Inst_Swap
	Inst_Drop
	Inst_Print
	Inst_PrintChar
	Inst_Debug
	Inst_Dump
	Inst_Label
	Inst_Start
	Inst_Com
	Inst_Alloc
	// MEM
	Inst_MemR8
	Inst_Var
	// Compilation only
	MemSet
)

func (ik InstKind) String() string {
	switch ik {
	case Inst_Dump:
		return "dump"
	case Inst_Push:
		return "push"
	case Inst_PushInt:
		return "pushi"
	case Inst_PushFloat:
		return "pushf"
	case Inst_PushUInt32:
		return "pushu"
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
	case Inst_Call:
		return "call"
	case Inst_JmpTrue:
		return "jmptrue"
	case Inst_Dup:
		return "dup"
	case Inst_Swap:
		return "swap"
	case Inst_Drop:
		return "drop"
	case Inst_Print:
		return "print"
	case Inst_PrintChar:
		return "printc"
	case Inst_Label:
		return "label"
	case Inst_Com:
		return "comment"
	case Inst_Ret:
		return "ret"
	case Inst_Start:
		return "__start:"
	case Inst_EqInt:
		return "eqi"
	case Inst_EqFloat:
		return "eqf"
	case Inst_Alloc:
		return "alloc"
	case Inst_Debug:
		return "debug"
	case Inst_MemR8:
		return "memr8"
	case Inst_Var:
		return "var"
	case MemSet:
		return "memset"
	default:
		fatal.Panic("InstKind unknown human representation of error: %d", ik)
	}
	return ""
}
