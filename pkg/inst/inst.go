package inst

import (
	"fmt"

	"github.com/fmarmol/vm/pkg/fatal"
	"github.com/fmarmol/vm/pkg/word"
)

type Inst struct {
	Kind    InstKind
	Operand word.Word // operand are `values` to be pushed on the stack
}

func NewInst(kind InstKind) func(word.Word) Inst {
	return func(value word.Word) Inst {
		return Inst{Kind: kind, Operand: value}
	}
}

var (
	// no operand
	Start = Inst{Kind: Inst_Start} // start is the entry point
	Add   = Inst{Kind: Inst_Add}   // add
	Sub   = Inst{Kind: Inst_Sub}   // substract
	Mul   = Inst{Kind: Inst_Mul}   // multiply
	Div   = Inst{Kind: Inst_Div}   // divide
	Print = Inst{Kind: Inst_Print} // print the value at the top of the stack
	Ret   = Inst{Kind: Inst_Ret}   // ret take the value at the top of the stack and assign ip to it. ret is used in functions to return to the caller next instruction
	Halt  = Inst{Kind: Inst_Halt}  // stop the vm
	Eq    = Inst{Kind: Inst_Eq}    // check if last 2 values are equal and but 1 or 0 at the top
	Drop  = Inst{Kind: Inst_Drop}  // remove value at the top of the stack
	Alloc = Inst{Kind: Inst_Alloc} // alloc the number of bytes value (uint32) at the top of the stack
	Dump  = Inst{Kind: Inst_Dump}  // dump the stack

	//operand
	PushInt    = NewInst(Inst_PushInt)    // push integer at the top of the stack
	PushFloat  = NewInst(Inst_PushFloat)  // push float at the top of the stack
	PushUInt32 = NewInst(Inst_PushUInt32) // push uint32 at the top of the stack
	Jmp        = NewInst(Inst_Jmp)        // Jmp at a position of the program
	JmpTrue    = NewInst(Inst_JmpTrue)    // Jump if top value of the stack != 0 at the position of the program
	Call       = NewInst(Inst_Call)       // call function
	Dup        = NewInst(Inst_Dup)        // Duplicate the value at the relative position in stack at the top of the stack
	Label      = NewInst(Inst_Label)      // label
	Swap       = NewInst(Inst_Swap)       //  swap the top of the stack with the relative position from sp
	EqInt      = NewInst(Inst_EqInt)      // compare the value with the top of the stack
	EqFloat    = NewInst(Inst_EqFloat)    // compare the value with the top of the stack

)

func (i Inst) String() string {
	switch i.Kind {
	// operand
	case Inst_PushInt, Inst_PushFloat, Inst_Jmp, Inst_JmpTrue, Inst_Dup, Inst_Label, Inst_Call, Inst_Swap, Inst_EqInt, Inst_EqFloat, Inst_PushUInt32:
		return fmt.Sprintf("%v %v", i.Kind, i.Operand)
	// no operand
	case Inst_Add, Inst_Halt, Inst_Sub, Inst_Mul, Inst_Div, Inst_Print, Inst_Drop, Inst_Ret, Inst_Start, Inst_Alloc, Inst_Dump:
		return fmt.Sprintf("%v", i.Kind)
	default:
		fatal.Panic("Inst unknown human representation of error: %v", i.Kind)
	}
	return ""
}
