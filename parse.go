package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fmarmol/regex"
	"github.com/fmarmol/tuple"
)

type Rule struct {
	kind    InstKind
	pattern string
	re      *regexp.Regexp
}

func loadRules() []*Rule {
	var rules = []*Rule{
		{kind: Inst_Com, pattern: `^(?P<com>//.*)`},
		{kind: Inst_Label, pattern: `^(?P<label>[[:word:]]+):`},
		{kind: Inst_JmpTrue, pattern: `^jmptrue\s+(?P<label>[[:word:]]+)`},
		{kind: Inst_Jmp, pattern: `^jmp\s+(?P<label>[[:word:]]+)`},
		{kind: Inst_PushInt, pattern: `^pushi\s+(?P<operand>\d+)`},
		{kind: Inst_PushFloat, pattern: `^pushf\s+(?P<operand>\d+)`},
		{kind: Inst_Dup, pattern: `^dup\s+(?P<operand>\d+)`},
		{kind: Inst_Add, pattern: `^(?P<inst>add)`},
		{kind: Inst_Sub, pattern: `^(?P<inst>sub)`},
		{kind: Inst_Print, pattern: `(?P<inst>^print)`},
		{kind: Inst_Halt, pattern: `^(?P<inst>halt)`},
	}
	for _, r := range rules {
		r.re = regexp.MustCompile(r.pattern)
	}
	return rules

}

func loadSourceCode(code string) *Program {
	labels := map[string]uint32{} // label: instruction position

	instsToResolve := []tuple.Tuple2[string, uint]{}

	var p Program

	lines := strings.Split(code, "\n")
	lines = func() (ret []string) {
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			ret = append(ret, line)
		}
		return
	}()

	var ip uint32

LINE:
	for _, line := range lines {
		var newInst Inst
		var found bool
		for _, rule := range loadRules() {
			groups := regex.FindGroups(rule.re, line)
			if len(groups) == 0 {
				continue
			}
			found = true
			switch rule.kind {
			case Inst_Com:
				ip++
				continue LINE
			case Inst_Label:
				label := groups.MustGet("label")

				_, ok := labels[label]
				if ok {
					Panic("label %v already defined", label)
				}
				labels[label] = ip
				newInst = Label(NewWord(ip, UInt32))
			case Inst_PushInt:
				op := groups.MustGetAsInt("operand")
				newInst = PushInt(NewWord(int64(op), Int64))
			case Inst_PushFloat:
				op := groups.MustGetAsFloat("operand")
				newInst = PushFloat(NewWord(op, Float64))
			case Inst_Dup:
				op := groups.MustGetAsInt("operand")
				newInst = Dup(NewWord(int64(op), Int64))
			case Inst_Jmp:
				label := groups.MustGet("label")

				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(len(p))))
				}
				newInst = Jmp(NewWord(addr, Int64))
			case Inst_JmpTrue:
				label := groups.MustGet("label")

				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(len(p))))
				}
				newInst = JmpTrue(NewWord(addr, Int64))
			case Inst_Add:
				newInst = Add
			case Inst_Sub:
				newInst = Sub
			case Inst_Halt:
				newInst = Halt
			case Inst_Print:
				newInst = Print
			default:
				Panic("Unkwon instruction line: %v", line)
			}
			p = append(p, newInst)
			ip++
			continue LINE
		}
		if !found {
			Panic("could not parse line: %v", line)
		}
	}

	for _, inst := range instsToResolve {
		res, ok := labels[inst.First]
		if !ok {
			Panic("Label %q is not defined", inst.First)
		}
		p[inst.Second].Operand = NewWord(res, UInt32)
	}

	for _, inst := range p {
		fmt.Println(inst)
	}

	return &p
}
