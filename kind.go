package main

type InstKind uint32

const (
	Inst_PushInt InstKind = iota
	Inst_PushFloat
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
	Inst_Label
	Inst_Start
	Inst_Com
	Inst_Count
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
	default:
		Panic("InstKind unknown human representation of error: %d", ik)
	}
	return ""
}
