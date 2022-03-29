package prog

import "github.com/fmarmol/vm/pkg/inst"

const PROG_CAPACITY = 1024

type Program [PROG_CAPACITY]inst.Inst

func (p *Program) Size() uint32 {
	return uint32(len(*p))
}

// func NewProgram(insts ...inst.Inst) *Program {
// 	ret := make([]inst.Inst, 0, len(insts))
// 	for _, inst := range insts {
// 		ret = append(ret, inst)
// 	}
// 	p := Program(ret)
// 	return &p
// }
