package prog

import (
	"regexp"
	"strings"

	"github.com/fmarmol/regex"
	"github.com/fmarmol/tuple"
	"github.com/fmarmol/vm/pkg/fatal"
	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/word"
)

type Rule struct {
	kind    inst.InstKind
	pattern string
	re      *regexp.Regexp
}

func loadRules() []*Rule {
	var rules = []*Rule{
		{kind: inst.Inst_Start, pattern: `^(?P<label>__start:)`},
		{kind: inst.Inst_Com, pattern: `^(?P<com>//.*)`},
		{kind: inst.Inst_Label, pattern: `^(?P<label>[[:word:]]+):`},
		{kind: inst.Inst_JmpTrue, pattern: `^jmptrue\s+(?P<label>[[:word:]]+)`},
		{kind: inst.Inst_Jmp, pattern: `^jmp\s+(?P<label>[[:word:]]+)`},
		{kind: inst.Inst_Call, pattern: `^call\s+(?P<label>[[:word:]]+)`},
		{kind: inst.Inst_PushUInt32, pattern: `^pushu\s+(?P<operand>\d+)`},
		{kind: inst.Inst_PushInt, pattern: `^pushi\s+(?P<operand>[-+]?\d+)`},
		{kind: inst.Inst_PushFloat, pattern: `^pushf\s+(?P<operand>[-+]?[0-9]+.[0-9]*)`},
		{kind: inst.Inst_Dup, pattern: `^dup\s+(?P<operand>\d+)`},
		{kind: inst.Inst_Swap, pattern: `^swap\s+(?P<operand>\d+)`},
		{kind: inst.Inst_EqFloat, pattern: `^eqf\s+(?P<operand>[-+]?[0-9]+.[0-9]*)`},
		{kind: inst.Inst_EqInt, pattern: `^eqi\s+(?P<operand>[-+]?\d+)`},
		{kind: inst.Inst_Drop, pattern: `^(?P<inst>drop)`},
		{kind: inst.Inst_Add, pattern: `^(?P<inst>add)`},
		{kind: inst.Inst_Sub, pattern: `^(?P<inst>sub)`},
		{kind: inst.Inst_Print, pattern: `^(?P<inst>^print)`},
		{kind: inst.Inst_Ret, pattern: `^(?P<inst>ret)`},
		{kind: inst.Inst_Halt, pattern: `^(?P<inst>halt)`},
		{kind: inst.Inst_Alloc, pattern: `^(?P<inst>alloc)`},
		{kind: inst.Inst_WriteStr, pattern: `^(?P<inst>write)\s+"(?P<operand>[[:word:]]+)"`},
		{kind: inst.Inst_Dump, pattern: `^(?P<inst>dump)`},
	}
	for _, r := range rules {
		r.re = regexp.MustCompile(r.pattern)
	}
	return rules

}

func LoadSourceCode(code string) *Program {
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
	var foundStart bool
	var foundStop bool

LINE:
	for _, line := range lines {
		var newInst inst.Inst
		var found bool
		for _, rule := range loadRules() {
			groups := regex.FindGroups(rule.re, line)
			if len(groups) == 0 {
				continue
			}
			found = true
			switch rule.kind {
			case inst.Inst_Com:
				ip++
				continue LINE
			case inst.Inst_Start:
				foundStart = true
				newInst = inst.Start
			case inst.Inst_Label:
				label := groups.MustGet("label")

				_, ok := labels[label]
				if ok {
					fatal.Panic("label %v already defined", label)
				}
				labels[label] = ip
				newInst = inst.Label(word.NewWord(ip, word.UInt32))
			case inst.Inst_PushInt:
				op := groups.MustGetAsInt("operand")
				newInst = inst.PushInt(word.NewWord(int64(op), word.Int64))
			case inst.Inst_PushFloat:
				op := groups.MustGetAsFloat("operand")
				newInst = inst.PushFloat(word.NewWord(op, word.Float64))
			case inst.Inst_PushUInt32:
				op := groups.MustGetAsInt("operand")
				newInst = inst.PushUInt32(word.NewWord(uint32(op), word.UInt32))
			case inst.Inst_EqInt:
				op := groups.MustGetAsInt("operand")
				newInst = inst.EqInt(word.NewWord(int64(op), word.Int64))
			case inst.Inst_EqFloat:
				op := groups.MustGetAsFloat("operand")
				newInst = inst.EqFloat(word.NewWord(op, word.Float64))
			case inst.Inst_Dup:
				op := groups.MustGetAsInt("operand")
				newInst = inst.Dup(word.NewWord(int64(op), word.UInt32))
			case inst.Inst_Jmp:
				label := groups.MustGet("label")

				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(len(p))))
				}
				newInst = inst.Jmp(word.NewWord(addr, word.UInt32))
			case inst.Inst_JmpTrue:
				label := groups.MustGet("label")

				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(len(p))))
				}
				newInst = inst.JmpTrue(word.NewWord(addr, word.UInt32))
			case inst.Inst_Call:
				label := groups.MustGet("label")

				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(len(p))))
				}
				newInst = inst.Call(word.NewWord(addr, word.UInt32))
			// case Inst_WriteStr:
			// 	op := groups.MustGet("operand")
			// 	newInst = WriteStr(NewWord(op, Str))
			case inst.Inst_Swap:
				idx := groups.MustGetAsInt("operand")
				newInst = inst.Swap(word.NewWord(uint32(idx), word.UInt32))
			case inst.Inst_Drop:
				newInst = inst.Drop
			case inst.Inst_Add:
				newInst = inst.Add
			case inst.Inst_Sub:
				newInst = inst.Sub
			case inst.Inst_Ret:
				newInst = inst.Ret
			case inst.Inst_Halt:
				newInst = inst.Halt
				foundStop = true
			case inst.Inst_Print:
				newInst = inst.Print
			case inst.Inst_Alloc:
				newInst = inst.Alloc
			case inst.Inst_Dump:
				newInst = inst.Dump
			default:
				fatal.Panic("Unkwon instruction line: %v", line)
			}
			if newInst.Kind == 0 {
				fatal.Panic("empty instruction")
			}
			p = append(p, newInst)
			ip++
			continue LINE
		}
		if !found {
			fatal.Panic("could not parse line: %v", line)
		}
	}
	if !foundStart {
		fatal.Panic("no entry point __start: found")
	}
	if !foundStop {
		fatal.Panic("no halt found")
	}

	for _, inst := range instsToResolve {
		res, ok := labels[inst.First]
		if !ok {
			fatal.Panic("Label %q is not defined", inst.First)
		}
		p[inst.Second].Operand = word.NewWord(res, word.UInt32)
	}

	return &p
}
