package mem

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestDump(t *testing.T) {
	var m Memory
	m.Write16(257, 0)
	m.Dump()
	assert.Equal(t, uint16(257), m.Read16(0))
}
