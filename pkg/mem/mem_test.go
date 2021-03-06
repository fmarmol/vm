package mem

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestDump(t *testing.T) {
	m := make(Memory, 10, 10)
	m.Write16(257, 0)
	m.Dump()
	assert.Equal(t, uint16(257), m.Read16(0))
}
