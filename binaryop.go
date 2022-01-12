package main

func binaryOp(arg1 int64, arg2 int64, ik InstKind) (ret int64, err Err) {
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
