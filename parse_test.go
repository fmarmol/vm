package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	code := `pushi 3
	pushi 2
	add
	`
	loadSourceCode(code)
}

func TestParseFloat(t *testing.T) {
	code := `pushf 3.14`
	p := loadSourceCode(code)
	assert.Len(t, *p, 1)
	assert.Equal(t, 3.14, (*p)[0].Operand.Float64())
}
