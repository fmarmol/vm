package rorre

import "github.com/fmarmol/vm/pkg/fatal"

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
	Err_AllocMem
)

func (e Err) Error() string { return e.String() }

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
	case Err_AllocMem:
		return "Error allocation memory"
	default:
		fatal.Panic("Err unknown human representation of error: %d", e)
	}
	return ""
}
