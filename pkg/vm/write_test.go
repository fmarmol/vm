package vm

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	v := InnerVM{}
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.BigEndian, v)
	assert.NoError(t, err)
}
