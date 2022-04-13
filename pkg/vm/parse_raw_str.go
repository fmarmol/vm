package vm

import (
	"fmt"
	"strings"
)

func parseRawStr(s string) (string, error) {
	if s[0] != '"' || s[len(s)-1] != '"' {
		return "", fmt.Errorf("could not parse into str: %q is not a string litteral", s)
	}
	v := strings.TrimFunc(s, func(r rune) bool {
		return r == '"'
	})
	return v, nil
}
