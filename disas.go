package main

import (
	"fmt"
)

func (p *Program) disas() (ret []string) {
	labels := map[Word]string{}
	indexToResolve := map[int64]Inst{}

	var ip int64

	fmt.Println("number of instructions:", len(*p))

	for _, inst := range *p {
		switch inst.Kind {
		case Inst_Label:
			labelName := fmt.Sprintf("__label_%d", ip)
			labels[inst.Operand] = labelName
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
		ip++
	}
	fmt.Println(ret)
	fmt.Println(labels)

	for index, inst := range indexToResolve {
		switch inst.Kind {
		case Inst_Jmp, Inst_JmpTrue:
			_, ok := labels[inst.Operand]
			if !ok {
				Panic("resolution of inst %v failed. Could not find labels at addr %v", inst, inst.Operand)
			}
			ret[index] = fmt.Sprintf("%v %v", inst.Kind, labels[inst.Operand])
		default:
			Panic("inst %v resolution is not implemented", inst)
		}
	}
	return
}
