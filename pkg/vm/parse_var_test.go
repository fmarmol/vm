package vm

import (
	"regexp"
	"testing"

	"github.com/fmarmol/regex"
	"github.com/stretchr/testify/assert"
)

func TestParseVar(t *testing.T) {
	{
		statement := `var x i64 = 3`

		vars := NewVars()

		re := regexp.MustCompile(VarDeclaration)
		groups := regex.FindGroups(re, statement)
		parseVar(vars, groups)

		assert.Len(t, vars.I64s, 1)
		assert.Contains(t, vars.I64s, "x")
		assert.Equal(t, int64(3), vars.I64s["x"].Value)
	}
	{
		statement := `var y f64 = 3.14`

		vars := NewVars()

		re := regexp.MustCompile(VarDeclaration)
		groups := regex.FindGroups(re, statement)
		parseVar(vars, groups)

		assert.Len(t, vars.F64s, 1)
		assert.Contains(t, vars.F64s, "y")
		assert.Equal(t, 3.14, vars.F64s["y"].Value)
	}
	{
		statement := `var x u32 = 3`

		vars := NewVars()

		re := regexp.MustCompile(VarDeclaration)
		groups := regex.FindGroups(re, statement)
		parseVar(vars, groups)

		assert.Len(t, vars.U32s, 1)
		assert.Contains(t, vars.U32s, "x")
		assert.Equal(t, uint32(3), vars.U32s["x"].Value)
	}
	{
		statement := `var msg str = "hello world"`

		vars := NewVars()

		re := regexp.MustCompile(VarDeclaration)
		groups := regex.FindGroups(re, statement)
		parseVar(vars, groups)

		assert.Len(t, vars.Strs, 1)
		assert.Contains(t, vars.Strs, "msg")
		assert.Equal(t, "hello world", vars.Strs["msg"].Value)
	}
}
