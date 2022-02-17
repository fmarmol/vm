package main

import "testing"

func TestParse(t *testing.T) {
	code := `pushi 3
	pushi 2
	add
	`
	loadSourceCode(code)
}
