package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/fmarmol/regex"
)

func Panic(s string, args ...interface{}) {
	err := fmt.Errorf(s, args...)
	panic(err)
}

type Rule struct {
	kind    InstKind
	pattern *regexp.Regexp
}

var pattern = regexp.MustCompile(`^(?P<inst>[[:alpha:]]+)( (?P<operand>\d+))?`)

func loadSourceCode(code string) *Program {

	ret := []Inst{}

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

	for _, line := range lines {
		groups := regex.FindGroups(pattern, line)
		if len(groups) == 0 {
			Panic("Cannot parse line: %v", line)
		}

		inst := groups.MustGet("inst")
		operand, ok := groups.Get("operand")

		var newInst Inst
		// no operand
		if ok {
			o, err := strconv.Atoi(operand)
			if err != nil {
				Panic("operand %q is not a number", operand)
			}
			op := int64(o)
			switch inst {
			case "push":
				newInst = Push(op)
			case "jmp":
				newInst = Jmp(op)
			case "jmptrue":
				newInst = JmpTrue(op)
			case "dup":
				newInst = Dup(op)
			default:
				Panic("unknow instruction %q", inst)
			}

		} else {
			switch inst {
			case "add":
				newInst = Add
			case "sub":
				newInst = Sub
			case "mul":
				newInst = Mul
			case "div":
				newInst = Div
			case "print":
				newInst = Print
			case "halt":
				newInst = Halt
			case "eq":
				newInst = Eq
			default:
				Panic("unknow instruction %v", inst)
			}
		}
		ret = append(ret, newInst)
	}
	return &ret
}
