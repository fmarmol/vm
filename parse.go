package main

import (
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
		{kind: Inst_Push, pattern: `^push\s+(?P<operand>\d+)`},
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
	labels := map[string]uint{}

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

LINE:
	for ip, line := range lines {
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
				continue LINE
			case Inst_Label:
				label := groups.MustGet("label")

				_, ok := labels[label]
				if ok {
					Panic("label %v already defined", label)
				}
				labels[label] = uint(ip)
				newInst = Label(int64(ip))
			case Inst_Push:
				op := groups.MustGetAsInt("operand")
				newInst = Push(int64(op))
			case Inst_Dup:
				op := groups.MustGetAsInt("operand")
				newInst = Dup(int64(op))
			case Inst_Jmp:
				label := groups.MustGet("label")

				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(len(p))))
				}
				newInst = Jmp(int64(addr))
			case Inst_JmpTrue:
				label := groups.MustGet("label")

				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(len(p))))
				}
				newInst = JmpTrue(int64(addr))
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
		p[inst.Second].Operand = int64(res)
	}

	return &p
}
