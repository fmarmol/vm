package main

import "testing"

func TestParse(t *testing.T) {
	code := `push 3
	push 2
	add
	`
	loadSourceCode(code)
}
