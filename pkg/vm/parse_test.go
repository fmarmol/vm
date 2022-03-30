package vm

import (
	"regexp"
	"testing"

	"github.com/fmarmol/regex"
	"gotest.tools/assert"
)

func TestParseSetM(t *testing.T) {
	s := `setmem 0 "hello world"`
	re := regexp.MustCompile(MemSetPattern)

	groups := regex.FindGroups(re, s)

	addr := groups.MustGetAsInt("addr")
	str := groups.MustGet("str")
	assert.Equal(t, 0, addr)
	assert.Equal(t, "hello world", str)

}
