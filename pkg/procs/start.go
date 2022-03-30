package procs

import (
	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/mem"
	"github.com/fmarmol/vm/pkg/word"
)

type VMer interface {
	IP() uint32
	ProgramSize() uint32
	StackPush(w word.Word) error
	SP() uint32
	StackCap() uint32
	// StackTop() word.Word
	Stop()                        // tell the vm to stop
	StackPop() (word.Word, error) // return the last elem of the stack and decrease sp
	// TODO: add error return
	StackPeek() word.Word                           // return the last elem without removing it
	StackPeekIndex(index uint32) (word.Word, error) // return the relative index to sp without removing it
	Swap(first, second uint32) error                // swap first and second index relative to sp (index >=1)
	Mem() *mem.Memory
	// Dup(index uint32) error                         // duplicate the index to relative to sp at the top of the stack
}

func Start(v VMer, _inst inst.Inst) error { return nil }
