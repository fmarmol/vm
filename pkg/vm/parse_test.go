package vm

import (
	"regexp"
	"testing"

	"github.com/fmarmol/regex"
	"github.com/stretchr/testify/assert"
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

func TestPush(t *testing.T) {
	re := regexp.MustCompile(PushPattern)
	type TestCase struct {
		name            string
		statement       string
		expectedOperand string
		expectedType    string
	}
	tcs := []TestCase{
		{
			name:            `push 1`,
			statement:       `push 1`,
			expectedOperand: "1",
			expectedType:    "",
		},
		{
			name:            `push 1[i64]`,
			statement:       `push 1[i64]`,
			expectedOperand: "1",
			expectedType:    "i64",
		},
		{
			name:            `push 1[u32]`,
			statement:       `push 1[u32]`,
			expectedOperand: "1",
			expectedType:    "u32",
		},
		{
			name:            `push 1[f64]`,
			statement:       `push 1[f64]`,
			expectedOperand: "1",
			expectedType:    "f64",
		},
		{
			name:            `push 3.14[f64]`,
			statement:       `push 3.14[f64]`,
			expectedOperand: "3.14",
			expectedType:    "f64",
		},
		{
			name:            `push 3.14`,
			statement:       `push 3.14`,
			expectedOperand: "3.14",
			expectedType:    "",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			groups := regex.FindGroups(re, tc.statement)
			op := groups.MustGet("operand")
			assert.Equal(t, tc.expectedOperand, op)
			kind, _ := groups.Get("type")
			assert.Equal(t, tc.expectedType, kind)
		})
	}
}

func TestParsePush(t *testing.T) {
	type TestCase[T ~int64 | uint32 | float64] struct {
		name          string
		statement     string
		expectedValue T
		expectedError bool
	}
	re := regexp.MustCompile(PushPattern)
	// int64
	tcsi := []TestCase[int64]{
		{
			name:          `i64 -> push 1`,
			statement:     `push 1`,
			expectedValue: 1,
			expectedError: false,
		},
		{
			name:          `push 1[i64]`,
			statement:     `push 1[i64]`,
			expectedValue: 1,
			expectedError: false,
		},
	}
	for _, tc := range tcsi {
		t.Run(tc.name, func(t *testing.T) {
			groups := regex.FindGroups(re, tc.statement)
			res, err := parsePush(tc.statement, groups)

			if !tc.expectedError {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.expectedValue, res.Operand.Int64())
			}
		})
	}
	// u32
	// TODO: check overflow
	tcsu := []TestCase[uint32]{
		{
			name:          `push 1[u32]`,
			statement:     `push 1[u32]`,
			expectedValue: 1,
			expectedError: false,
		},
	}
	for _, tc := range tcsu {
		t.Run(tc.name, func(t *testing.T) {
			groups := regex.FindGroups(re, tc.statement)
			res, err := parsePush(tc.statement, groups)

			if !tc.expectedError {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.expectedValue, res.Operand.UInt32())
			}
		})
	}
	//f64
	tcsf := []TestCase[float64]{
		{
			name:          `push 3.14`,
			statement:     `push 3.14`,
			expectedValue: 1,
			expectedError: false,
		},
		{
			name:          `push 3.14[f64]`,
			statement:     `push 3.14[f64]`,
			expectedValue: 1,
			expectedError: false,
		},
	}
	for _, tc := range tcsf {
		t.Run(tc.name, func(t *testing.T) {
			groups := regex.FindGroups(re, tc.statement)
			res, err := parsePush(tc.statement, groups)

			if !tc.expectedError {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.expectedValue, res.Operand.Float64())
			}
		})
	}
}
