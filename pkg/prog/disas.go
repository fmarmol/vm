package prog

import (
	"fmt"

	"github.com/fmarmol/vm/pkg/fatal"
	"github.com/fmarmol/vm/pkg/inst"
)

func (p *Program) Disas() (ret []string) {
	labels := map[uint32]string{}
	indexToResolve := map[uint32]inst.Inst{}

	var ip uint32

	for _, _inst := range *p {
		switch _inst.Kind {
		case inst.Inst_Label:
			labelName := fmt.Sprintf("__label_%d", ip)
			labels[_inst.Operand.UInt32()] = labelName
			ret = append(ret, labelName+":")
		case inst.Inst_Jmp, inst.Inst_JmpTrue, inst.Inst_Call:
			//check if addres is saved in labels:
			_, ok := labels[_inst.Operand.UInt32()]
			if !ok { // need to resolve later
				indexToResolve[ip] = _inst
				ret = append(ret, "")
			} else {
				ret = append(ret, fmt.Sprintf("%v %s", _inst.Kind, labels[_inst.Operand.UInt32()]))
			}
		default:
			ret = append(ret, fmt.Sprintf("%v", _inst))
		}
		ip++
	}
	for index, _inst := range indexToResolve {
		switch _inst.Kind {
		case inst.Inst_Jmp, inst.Inst_JmpTrue:
			_, ok := labels[_inst.Operand.UInt32()]
			if !ok {
				fatal.Panic("resolution of inst %v failed. Could not find labels at addr %v", _inst, _inst.Operand)
			}
			ret[index] = fmt.Sprintf("%v %v", _inst.Kind, labels[_inst.Operand.UInt32()])
		default:
			fatal.Panic("inst %v resolution is not implemented", _inst)
		}
	}
	return
}
