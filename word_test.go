package main

import (
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
	w := NewWord(a, Int64)
	assert.Equal(t, a, w.UInt32())
}
