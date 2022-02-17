package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var a float32 = -1.0
	fmt.Printf("%016b\n", a)

	ref := *(*int32)(unsafe.Pointer(&a))
	tot := unsafe.Sizeof(ref) * 8

	for i := int(tot - 1); i >= 0; i-- {
		res := (ref >> i) & 1
		fmt.Printf("%b", res)
	}
	fmt.Println("")
	fmt.Printf("ref: %d\n", ref)

	restore := *(*float32)(unsafe.Pointer(&ref))
	fmt.Printf("restore: %f\n", restore)
	// fmt.Printf("%016b\n", a)
}
