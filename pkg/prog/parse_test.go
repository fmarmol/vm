package prog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	code := `
	__start:
	pushi 3
	pushi 2
	add
	halt
	`
	LoadSourceCode(code)
}

func TestParseFloat(t *testing.T) {
	code := `
	__start:
	pushf 3.14
	halt
	`
	p := LoadSourceCode(code)
	assert.Len(t, *p, 3)
	assert.Equal(t, 3.14, (*p)[1].Operand.Float64())
}
