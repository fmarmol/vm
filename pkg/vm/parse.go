package vm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/fmarmol/regex"
	"github.com/fmarmol/tuple"
	"github.com/fmarmol/vm/pkg/fatal"
	"github.com/fmarmol/vm/pkg/inst"
	"github.com/fmarmol/vm/pkg/mem"
	"github.com/fmarmol/vm/pkg/prog"
	"github.com/fmarmol/vm/pkg/word"
)

const MemSetPattern = `^setmem\s+(?P<addr>\d+)\s+"(?P<str>[^"]*)"`
const PushPattern = `^push\s+(?P<operand>[^\[]+)(\[(?P<type>(i64|u32|f64))\])?`

func parsePush(statement string, groups regex.Groups) (inst.Inst, error) {
	operand := groups.MustGet("operand")
	kind, ok := groups.Get("type")

	if !ok {
		{
			res, err := strconv.ParseInt(operand, 10, 64)
			if err == nil {
				return inst.PushInt(word.NewI64(res)), nil
			}
		}
		{
			res, err := strconv.ParseFloat(operand, 64)
			if err == nil {
				return inst.PushFloat(word.NewF64(res)), nil
			}

		}
	} else {
		switch kind {
		case "i64":
			res, err := strconv.ParseInt(operand, 10, 64)
			if err == nil {
				return inst.PushInt(word.NewI64(res)), nil
			} else {
				return inst.Inst{}, fmt.Errorf("could not convert [%v] into i64: %v", operand, err)
			}
		case "u32":
			res, err := strconv.ParseUint(operand, 10, 32)
			if err == nil {
				return inst.PushUInt32(word.NewU32(uint32(res))), nil
			} else {
				return inst.Inst{}, fmt.Errorf("could not convert [%v] into u32: %v", operand, err)
			}
		case "f64":
			res, err := strconv.ParseFloat(operand, 64)
			if err == nil {
				return inst.PushFloat(word.NewF64(res)), nil
			} else {
				return inst.Inst{}, fmt.Errorf("could not convert [%v] into f64: %v", operand, err)
			}
		default:
			return inst.Inst{}, fmt.Errorf("could not parse push because unknown type: %v", kind)

		}
	}
	return inst.Inst{}, fmt.Errorf("could not parse push statement: %q", statement)
}

func loadRules() []*Rule {
	var rules = []*Rule{
		{kind: inst.Inst_Start, pattern: `^(?P<label>__start:)`},
		{kind: inst.Inst_Com, pattern: `^(?P<com>//.*)`},
		{kind: inst.Inst_Label, pattern: `^(?P<label>[[:word:]]+):`},
		{kind: inst.Inst_JmpTrue, pattern: `^jmptrue\s+(?P<label>[[:word:]]+)`},
		{kind: inst.Inst_Jmp, pattern: `^jmp\s+(?P<label>[[:word:]]+)`},
		{kind: inst.Inst_Call, pattern: `^call\s+(?P<label>[[:word:]]+)`},
		{kind: inst.Inst_Push, pattern: PushPattern}, // default value is i64 and f64
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
		{kind: inst.Inst_Print, pattern: `^(?P<inst>print)`},
		{kind: inst.Inst_PrintChar, pattern: `^(?P<inst>printc)`},
		{kind: inst.Inst_Debug, pattern: `^(?P<inst>debug)`},
		{kind: inst.Inst_Ret, pattern: `^(?P<inst>ret)`},
		{kind: inst.Inst_Halt, pattern: `^(?P<inst>halt)`},
		{kind: inst.Inst_Alloc, pattern: `^(?P<inst>alloc)`},
		{kind: inst.Inst_Dump, pattern: `^(?P<inst>dump)`},
		{kind: inst.MemSet, pattern: MemSetPattern},
		{kind: inst.Inst_MemR8, pattern: `^(?P<inst>memr8)`},
		{kind: inst.Inst_DisplayStr, pattern: `^x\\str\s+$(?P<variable>[[:word:]]+)`},
	}
	for _, r := range rules {
		r.re = regexp.MustCompile(r.pattern)
	}
	return rules
}

type Variable[T any] struct {
	ptr   uint32
	size  uint64
	value T
}

func LoadSourceCode(code string) InnerVM {
	labels := map[string]uint32{} // label: instruction position
	variable := map[string]Variable{}

	instsToResolve := []tuple.Tuple2[string, uint]{}

	var p prog.Program

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

	var m mem.Memory

	var ip uint32
	var foundStart bool
	var foundStop bool
	// var i int

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
				newInst = inst.Label(word.NewU32(ip))
			case inst.Inst_Push:
				_inst, err := parsePush(line, groups) // TODO: This function's signature is weird
				if err != nil {
					panic(err)
				}
				newInst = _inst
			case inst.Inst_EqInt:
				op := groups.MustGetAsInt("operand")
				newInst = inst.EqInt(word.NewI64(int64(op)))
			case inst.Inst_EqFloat:
				op := groups.MustGetAsFloat("operand")
				newInst = inst.EqFloat(word.NewF64(op))
			case inst.Inst_Dup:
				op := groups.MustGetAsInt("operand")
				newInst = inst.Dup(word.NewU32(uint32(op)))
			case inst.Inst_Jmp:
				label := groups.MustGet("label")
				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(ip)))
				}
				newInst = inst.Jmp(word.NewU32(addr))
			case inst.Inst_JmpTrue:
				label := groups.MustGet("label")

				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(ip)))
				}
				newInst = inst.JmpTrue(word.NewU32(addr))
			case inst.Inst_Call:
				label := groups.MustGet("label")

				addr, ok := labels[label]
				if !ok {
					instsToResolve = append(instsToResolve, tuple.NewTuple2(label, uint(ip)))
				}
				newInst = inst.Call(word.NewU32(addr))
			case inst.Inst_Swap:
				idx := groups.MustGetAsInt("operand")
				newInst = inst.Swap(word.NewU32(uint32(idx)))
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
			case inst.Inst_Debug:
				newInst = inst.Debug
			case inst.Inst_Alloc:
				newInst = inst.Alloc
			case inst.Inst_Dump:
				newInst = inst.Dump
			case inst.MemSet:
				// TODO: MANAGE SPECIAL CARACTERS
				addr := groups.MustGetAsInt("addr")
				str := groups.MustGet("str")
				for index := 0; index < len(str); index++ {
					if addr+index > len(m)-1 {
						cop
					}
					m.Write8(str[index], uint32(addr+index))
				}
				ip++
				continue LINE
			case inst.Inst_MemR8:
				newInst = inst.MemR8

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
		p[inst.Second].Operand = word.NewU32(res)
	}
	return InnerVM{Program: p, Memory: m}
}
