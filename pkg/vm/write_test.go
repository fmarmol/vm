package vm

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	v := NewVM(InnerVM{})
	buf := bytes.NewBuffer(nil)
	err := v.WriteToFile(buf)
	assert.NoError(t, err)
}
