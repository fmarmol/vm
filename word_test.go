package main

import (
	"encoding/binary"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewWordInt64(t *testing.T) {

	var a int64 = 1
	w := NewWord(a, Int64)
	assert.Equal(t, a, w.Int64())
}

func TestNewWordFloat64(t *testing.T) {
	var a float64 = 3.14
	w := NewWord(a, Float64)
	assert.Equal(t, a, w.Float64())
}
func TestNewWordUInt32(t *testing.T) {
	var a uint32 = 2
	w := NewWord(a, UInt32)
	assert.Equal(t, a, w.UInt32())
}

func TestUint64SizeOf(t *testing.T) {
	assert.Equal(t, 8, binary.Size(uint64(0)))
}

func TestWordKindSizeOf(t *testing.T) {
	assert.Equal(t, 1, binary.Size(UInt32))
}

func TestSizeOfWord(t *testing.T) {
	w := Word{Kind: UInt32, Value: 1}
	assert.Equal(t, int((64+8)/8), binary.Size(w))
}
