package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRawStr(t *testing.T) {
	s := `"hello"`

	res, err := parseRawStr(s)
	assert.NoError(t, err)
	assert.Equal(t, "hello", res)
}
