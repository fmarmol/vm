package vm

import (
	"fmt"
	"strconv"

	"github.com/fmarmol/regex"
	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/word"
)

func parsePush(statement string, groups regex.Groups) (inst.Inst, error) {
	operand := groups.MustGet("operand")
	kind, ok := groups.Get("type")

	if !ok {
		{
			res, err := strconv.ParseInt(operand, 10, 64)
			if err == nil {
				return inst.PushInt(word.NewI64(res)), nil
			}
		}
		{
			res, err := strconv.ParseFloat(operand, 64)
			if err == nil {
				return inst.PushFloat(word.NewF64(res)), nil
			}

		}
	} else {
		switch kind {
		case "i64":
			res, err := strconv.ParseInt(operand, 10, 64)
			if err == nil {
				return inst.PushInt(word.NewI64(res)), nil
			} else {
				return inst.Inst{}, fmt.Errorf("could not convert [%v] into i64: %v", operand, err)
			}
		case "u32":
			res, err := strconv.ParseUint(operand, 10, 32)
			if err == nil {
				return inst.PushUInt32(word.NewU32(uint32(res))), nil
			} else {
				return inst.Inst{}, fmt.Errorf("could not convert [%v] into u32: %v", operand, err)
			}
		case "f64":
			res, err := strconv.ParseFloat(operand, 64)
			if err == nil {
				return inst.PushFloat(word.NewF64(res)), nil
			} else {
				return inst.Inst{}, fmt.Errorf("could not convert [%v] into f64: %v", operand, err)
			}
		default:
			return inst.Inst{}, fmt.Errorf("could not parse push because unknown type: %v", kind)

		}
	}
	return inst.Inst{}, fmt.Errorf("could not parse push statement: %q", statement)
}
