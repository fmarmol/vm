package main

import (
	"fmt"
)

func (p *Program) disas() (ret []string) {
	labels := map[uint32]string{}
	indexToResolve := map[uint32]Inst{}

	var ip uint32

	for _, inst := range *p {
		switch inst.Kind {
		case Inst_Label:
			labelName := fmt.Sprintf("__label_%d", ip)
			labels[inst.Operand.UInt32()] = labelName
			ret = append(ret, labelName+":")
		case Inst_Jmp, Inst_JmpTrue, Inst_Call:
			//check if addres is saved in labels:
			_, ok := labels[inst.Operand.UInt32()]
			if !ok { // need to resolve later
				indexToResolve[ip] = inst
				ret = append(ret, "")
			} else {
				ret = append(ret, fmt.Sprintf("%v %s", inst.Kind, labels[inst.Operand.UInt32()]))
			}
		default:
			ret = append(ret, fmt.Sprintf("%v", inst))
		}
		ip++
	}
	for index, inst := range indexToResolve {
		switch inst.Kind {
		case Inst_Jmp, Inst_JmpTrue:
			_, ok := labels[inst.Operand.UInt32()]
			if !ok {
				Panic("resolution of inst %v failed. Could not find labels at addr %v", inst, inst.Operand)
			}
			ret[index] = fmt.Sprintf("%v %v", inst.Kind, labels[inst.Operand.UInt32()])
		default:
			Panic("inst %v resolution is not implemented", inst)
		}
	}
	return
}
