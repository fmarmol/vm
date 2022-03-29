package vm

// func (v *VM) executeInst(_inst inst.Inst) (f func(v *VM, _inst inst.Inst) error, _err rorre.Err) {
// 	switch _inst.Kind {
// 	// case inst.Inst_Alloc:
// 	// 	nbytes := v.stack[v.sp-1].UInt32()
// 	// 	v.stack[v.sp-1] = word.Word{}
// 	// 	v.sp--
// 	// 	data, _err2 := AllocMem(nbytes)
// 	// 	if _err2 != nil {
// 	// 		_err = Err_AllocMem
// 	// 		return
// 	// 	}
// 	// 	addr := uintptr(unsafe.Pointer(&data[0]))
// 	// 	v.stack[v.sp] = NewWord(addr, Ptr)
// 	// 	v.sp++
// 	// 	v.ip++
// 	case inst.Inst_Start:
// 		v.ip++
// 	case inst.Inst_PushInt, inst.Inst_PushFloat, inst.Inst_PushUInt32:
// 		if v.sp >= v.maxSize {
// 			_err = rorre.Err_Overflow
// 		} else {
// 			v.stack[v.sp] = _inst.Operand
// 			v.sp++
// 			v.ip++
// 		}
// 	// DEBUG INSTRUCTIONS
// 	case inst.Inst_Dump:
// 		v.dump()
// 		v.ip++
// 	case inst.Inst_EqInt:
// 		a := v.stack[v.sp-1]
// 		op := _inst.Operand
// 		if a.Kind != op.Kind {
// 			fmt.Fprintf(os.Stderr, "incompatible types: tried comparison between %v and %v\n", a.Kind, op.Kind)
// 			_err = rorre.Err_WrongTypeOperation
// 			return
// 		}
// 		if a.Int64() != op.Int64() {
// 			fatal.Panic("invalid assertion top[%v] != eq[%v]", a.Int64(), op.Int64())
// 		}
// 		v.ip++
// 	case inst.Inst_EqFloat:
// 		top := v.stackTop()
// 		op := _inst.Operand
// 		if top.Kind != op.Kind {
// 			fmt.Fprintf(os.Stderr, "incompatible types: tried comparison between %v and %v\n", top.Kind, op.Kind)
// 			_err = rorre.Err_WrongTypeOperation
// 			return
// 		}
// 		if top.Float64() != op.Float64() {
// 			fatal.Panic("invalid assertion top[%v] != eq[%v]", top.Float64(), op.Float64())
// 		}
// 		v.ip++
// 	// END DEBUG INSTRUCTIONS
// 	case inst.Inst_Add, inst.Inst_Sub, inst.Inst_Mul, inst.Inst_Div, inst.Inst_Eq:
// 		if len(v.stack) < 2 {
// 			_err = rorre.Err_Underflow
// 			return
// 		}
// 		a := v.stack[v.sp-2]
// 		b := v.stack[v.sp-1]
// 		if a.Kind != b.Kind {
// 			fmt.Fprintf(os.Stderr, "incompatible types: tried to binary operation between %v and %v\n", a.Kind, b.Kind)
// 			_err = rorre.Err_WrongTypeOperation
// 			return
// 		}
// 		var result word.Word
// 		switch a.Kind {
// 		case word.Int64:
// 			res, _err2 := binaryOp(a.Int64(), b.Int64(), _inst.Kind)
// 			if _err2 != rorre.OK {
// 				_err = _err2
// 				return
// 			}
// 			result = word.NewWord(res, word.Int64)
// 		case word.Float64:
// 			res, _err2 := binaryOp(a.Float64(), b.Float64(), _inst.Kind)
// 			if _err2 != rorre.OK {
// 				_err = _err2
// 				return
// 			}
// 			result = word.NewWord(res, word.Float64)
// 		}
// 		v.stack[v.sp-2] = result
// 		v.stack[v.sp-1] = word.Word{}
// 		v.sp--
// 		v.ip++
// 	case inst.Inst_Swap:
// 		pos_top := v.sp - 1
// 		pos_sec := v.sp - (_inst.Operand.UInt32())
// 		v.stack[pos_sec], v.stack[pos_top] = v.stack[pos_top], v.stack[pos_sec]
// 		v.ip++
// 	case inst.Inst_Drop:
// 		v.stack[v.sp-1] = word.Word{}
// 		v.sp--
// 		v.ip++
// 	case inst.Inst_Halt:
// 		v.stop = true
// 	case inst.Inst_Ret:
// 		v.ip = v.stackTop().UInt32()
// 		v.stack[v.sp-1] = word.Word{}
// 		v.sp--
// 	case inst.Inst_Call:
// 		if _inst.Operand.UInt32() < 0 || _inst.Operand.UInt32() >= v.program.Size() {
// 			_err = rorre.Err_OutOfIndexInstruction
// 		} else {
// 			v.stack[v.sp] = word.NewWord(v.ip+1, word.UInt32)
// 			v.sp++
// 			v.ip = _inst.Operand.UInt32()
// 		}

// 	case inst.Inst_Jmp:
// 		if _inst.Operand.UInt32() < 0 || _inst.Operand.UInt32() >= v.program.Size() {
// 			_err = rorre.Err_OutOfIndexInstruction
// 		} else {
// 			v.ip = _inst.Operand.UInt32()
// 		}
// 	case inst.Inst_JmpTrue:
// 		if _inst.Operand.UInt32() < 0 || _inst.Operand.UInt32() >= v.program.Size() {
// 			_err = rorre.Err_OutOfIndexInstruction
// 		} else if !v.stack[v.sp-1].IsZero() {
// 			v.ip = _inst.Operand.UInt32()
// 		} else {
// 			v.ip++
// 		}
// 	case inst.Inst_Dup: // duplicate relative to sp
// 		if _inst.Operand.UInt32() <= 0 {
// 			_err = rorre.Err_OutOfIndexInstruction
// 		} else {
// 			v.stack[v.sp] = v.stack[v.sp-_inst.Operand.UInt32()]
// 			v.sp++
// 			v.ip++
// 		}
// 	case inst.Inst_Print:
// 		fmt.Println("->", v.stack[v.sp-1])
// 		v.ip++
// 	case inst.Inst_Label:
// 		v.ip++
// 	default:
// 		_err = rorre.Err_IllegalInstruction
// 	}
// 	return
// }
