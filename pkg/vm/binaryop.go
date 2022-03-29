package vm

import (
	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/rorre"
)

type Number interface {
	~uint32 | ~int64 | ~float64
}

func binaryOp[T Number](arg1 T, arg2 T, ik inst.InstKind) (ret T, err rorre.Err) {
	switch ik {
	case inst.Inst_Add:
		ret = arg1 + arg2
	case inst.Inst_Sub:
		ret = arg1 - arg2
	case inst.Inst_Mul:
		ret = arg1 * arg2
	case inst.Inst_Div:
		if arg2 == 0 {
			err = rorre.Err_DivisionByZero
		} else {
			ret = arg1 / arg2
		}
	case inst.Inst_Eq:
		if arg1 == arg2 {
			ret = 1
		} else {
			ret = 0
		}
	default:
		panic("unknown binaryOp")
	}
	return
}
