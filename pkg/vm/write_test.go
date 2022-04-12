package vm

import (
	"bytes"
	"testing"

	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/prog"
	"github.com/stretchr/testify/assert"
)

func TestWriteAndLoad(t *testing.T) {
	v := NewVM(InnerVM{
		Memory:  []byte{1, 2, 3},
		Program: prog.Program{inst.Add, inst.Mul},
	})
	buf := bytes.NewBuffer(nil)
	err := v.Write(buf)
	assert.NoError(t, err)

	nv, err := Load(buf)
	assert.NoError(t, err)
	assert.Equal(t, v.MetaInnerVM, nv.MetaInnerVM)
	assert.Equal(t, v.Memory, nv.Memory)
	assert.Equal(t, v.Program, nv.Program)
}
