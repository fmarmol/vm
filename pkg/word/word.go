package word

import (
	"fmt"
	"unsafe"
)

type WordKind int8

const (
	Int64 WordKind = iota
	Float64
	UInt32
	Ptr
)

func (w WordKind) String() string {
	switch w {
	case Int64:
		return "int64"
	case Float64:
		return "float64"
	case UInt32:
		return "uint32"
	case Ptr:
		return "ptr"
	default:
		panic(fmt.Errorf("unknown human representation fo WordKind %d", w))
	}
}

type Word struct {
	Kind  WordKind
	Value uint64
}

func (w Word) String() string {
	switch w.Kind {
	case Int64:
		return fmt.Sprintf("%d", w.Int64())
	case Float64:
		return fmt.Sprintf("%f", w.Float64())
	case UInt32:
		return fmt.Sprintf("%d", w.UInt32())
	case Ptr:
		return fmt.Sprintf("0x%02X", w.Ptr())
	default:
		panic("unknown human representation of word")
	}
}

func (w Word) IsZero() bool {
	return w.Value == 0
}

func newWord[T ~int64 | ~float64 | ~uint32 | uintptr](i T, kind WordKind) Word {
	w := Word{
		Kind:  kind,
		Value: *(*uint64)(unsafe.Pointer(&i)),
	}
	return w
}

func NewPtr(i uintptr) Word { return newWord(i, Ptr) }
func NewU32(i uint32) Word  { return newWord(i, UInt32) }
func NewF64(i float64) Word { return newWord(i, Float64) }
func NewI64(i int64) Word   { return newWord(i, Int64) }

func (w Word) Float64() float64 {
	if w.Kind != Float64 {
		panic(fmt.Errorf("cannot convert word value into float64 should be %v", w.Kind))
	}
	return *(*float64)(unsafe.Pointer(&w.Value))
}

func (w Word) Int64() int64 {
	if w.Kind != Int64 {
		panic(fmt.Errorf("cannot convert word value into int64 should be %v", w.Kind))
	}
	return *(*int64)(unsafe.Pointer(&w.Value))
}

func (w Word) UInt32() uint32 {
	if w.Kind != UInt32 {
		panic(fmt.Errorf("cannot convert word value into uint32 should be %v", w.Kind))
	}
	return *(*uint32)(unsafe.Pointer(&w.Value))
}

func (w Word) Ptr() uintptr {
	if w.Kind != Ptr {
		panic(fmt.Errorf("cannot convert word value into ptr should be %v", w.Kind))
	}
	return uintptr(unsafe.Pointer(&w.Value))
}
