package vm

import (
	"fmt"
	"strconv"

	"github.com/fmarmol/regex"
)

type VarI interface {
	~int64 | uint32 | float64 | string
}

type Var[T VarI] struct {
	Name  string
	Value T
	Ptr   uint32
}

type Vars struct {
	I64s map[string]Var[int64]
	U32s map[string]Var[uint32]
	F64s map[string]Var[float64]
	Strs map[string]Var[string]
}

func NewVars() *Vars {
	return &Vars{
		I64s: make(map[string]Var[int64]),
		U32s: make(map[string]Var[uint32]),
		F64s: make(map[string]Var[float64]),
		Strs: make(map[string]Var[string]),
	}
}

func parseVar(vars *Vars, groups regex.Groups) error {
	id := groups.MustGet("identifier")
	_type := groups.MustGet("type")
	value := groups.MustGet("value")

	switch _type {
	case "i64":
		res, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			vars.I64s[id] = Var[int64]{Value: res}
		} else {
			return fmt.Errorf("could not convert [%v] into i64: %v", value, err)
		}
	case "u32":
		res, err := strconv.ParseUint(value, 10, 32)
		if err == nil {
			vars.U32s[id] = Var[uint32]{Value: uint32(res)}
		} else {
			return fmt.Errorf("could not convert [%v] into u32: %v", value, err)
		}
	case "f64":
		res, err := strconv.ParseFloat(value, 64)
		if err == nil {
			vars.F64s[id] = Var[float64]{Value: res}
		} else {
			return fmt.Errorf("could not convert [%v] into f64: %v", value, err)
		}
	case "str":
		res, err := parseRawStr(value)
		if err == nil {
			vars.Strs[id] = Var[string]{Value: res}
		} else {
			return fmt.Errorf("could not convert [%v] into str: %v", value, err)
		}
	default:
		return fmt.Errorf("could not parse push because unknown type: %v", _type)
	}
	return nil
}
