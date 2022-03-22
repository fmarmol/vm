package fatal

import (
	"fmt"
	"os"
)

func Panic(s string, args ...interface{}) {
	err := fmt.Errorf(s, args...)
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}
