package main

type Number interface {
	~uint32 | ~int64 | ~float64
}

func binaryOp[T Number](arg1 T, arg2 T, ik InstKind) (ret T, err Err) {
	switch ik {
	case Inst_Add:
		ret = arg1 + arg2
	case Inst_Sub:
		ret = arg1 - arg2
	case Inst_Mul:
		ret = arg1 * arg2
	case Inst_Div:
		if arg2 == 0 {
			err = Err_DivisionByZero
		} else {
			ret = arg1 / arg2
		}
	case Inst_Eq:
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
