package vm

import (
	"errors"
	"regexp"

	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/mem"
	"github.com/fmarmol/vm/pkg/prog"
	"github.com/fmarmol/vm/pkg/rorre"
	"github.com/fmarmol/vm/pkg/word"
)

const STACK_CAPACITY = 1024 // 1024 Bytes should be enough for every one

type VM struct {
	Stack [STACK_CAPACITY]word.Word
	bp    uint32 // stack base pointer
	sp    uint32 // stack pointer
	ip    uint32 // instruction pointer
	stop  bool
	InnerVM
}

type InnerVM struct {
	Memory  mem.Memory
	Program prog.Program
}

type Rule struct {
	kind    inst.InstKind
	pattern string
	re      *regexp.Regexp
}

func (v *VM) Mem() *mem.Memory    { return &v.Memory }
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
