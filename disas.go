package main

import (
	"fmt"
)

func (p *Program) disas() (ret []string) {
	labels := map[int64]string{}
	indexToResolve := map[int64]Inst{}

	for ip, inst := range *p {
		switch inst.Kind {
		case Inst_Label:
			labelName := fmt.Sprintf("__label_%d", ip)
			labels[int64(ip)] = labelName
			ret = append(ret, labelName+":")
		case Inst_Jmp, Inst_JmpTrue:
			//check if addres is saved in labels:
			_, ok := labels[inst.Operand]
			if !ok { // need to resolve later
				indexToResolve[int64(ip)] = inst
				ret = append(ret, "")
			} else {
				ret = append(ret, fmt.Sprintf("%v %s", inst.Kind, labels[inst.Operand]))
			}
		default:
			ret = append(ret, fmt.Sprintf("%v", inst))
		}
	}

	for index, inst := range indexToResolve {
		switch inst.Kind {
		case Inst_Jmp, Inst_JmpTrue:
			_, ok := labels[inst.Operand]
			if !ok {
				Panic("resolution of inst %v failed. Could not find labels at addr %v", inst, inst.Operand)
			}
			ret[index] = fmt.Sprintf("%v %s", inst.Kind, labels[inst.Operand])
		default:
			Panic("inst %v resolution is not implemented", inst)
		}
	}
	return
}