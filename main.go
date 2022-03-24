package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fmarmol/basename/pkg/basename"
	"github.com/fmarmol/vm/pkg/fatal"
	"github.com/fmarmol/vm/pkg/prog"
	"github.com/fmarmol/vm/pkg/vm"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	NULL             uint32 = 0
	PROGRAM_CAPACITY        = 512
	STACK_CAPACITY          = 100
)

// func NewArray(size int, value byte) []byte {
// 	ret := make([]byte, 0, size)
// 	for i := 0; i < size; i++ {
// 		ret = append(ret, value)
// 	}
// 	return ret
// }

// func (v *VM) freeMemSpace(ptr uint32, size uint32) {
// 	copy(v.mem[int(ptr):int(ptr+size)], NewArray(int(size), 0))
// }

// func (v *VM) findMemSpace(size uint32) (uint32, Err) { // return the index on the memory able to allocate size bytes
// LOOP:
// 	for i := range v.mem {
// 		for j := range v.mem[i : i+int(size)] {
// 			if v.mem[i+j] != 0 {
// 				continue LOOP
// 			}
// 		}
// 		return v.memRegister.Register(uint32(i), size)
// 	}
// 	return 0, Err_SpaceNotFound
// }

var (
	app    = kingpin.New("vm", "vm main command")
	comp   = app.Command("compile", "compile a .evm file").Alias("c")
	source = comp.Arg("source", "source file").String()
	output = comp.Flag("output", "output file .vm").Short('o').String()

	run       = app.Command("run", "run vm file").Alias("r")
	sourceRun = run.Arg("source", "source file .vm").String()
	maxStep   = run.Flag("max_step", "max exection steps allowed").Default("300").Uint()

	debug        = app.Command("debug", "run vm file").Alias("d")
	sourceDebug  = debug.Arg("source", "source file .vm").String()
	maxStepDebug = debug.Flag("max_step", "max exection steps allowed").Default("300").Uint()

	disas       = app.Command("disas", "disassemble a program .vm")
	sourceDisas = disas.Arg("source", "source file .vm").String()
	outputDisas = disas.Flag("output", "output file .vm.disas").Short('o').String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case comp.FullCommand():
		fi := basename.ParseFile(*source)
		code, err := ioutil.ReadFile(*source)
		if err != nil {
			fatal.Panic("could not read file: %v", err)
		}
		p := prog.LoadSourceCode(string(code))
		v := vm.NewVM(PROGRAM_CAPACITY, p)
		path := filepath.Join(fi.Dir, fi.Basename) + ".vm"
		err = v.WriteToFile(path)
		if err != nil {
			panic(err)
		}
	case run.FullCommand():
		p, err := prog.LoadProgram(*sourceRun)
		if err != nil {
			panic(err)
		}
		v := vm.NewVM(PROGRAM_CAPACITY, p)
		v.Execute(*maxStep)
	case debug.FullCommand():
		p, err := prog.LoadProgram(*sourceDebug)
		if err != nil {
			panic(err)
		}
		v := vm.NewVM(PROGRAM_CAPACITY, p)
		v.ExecuteWithDebug(*maxStepDebug)
	case disas.FullCommand():
		p, err := prog.LoadProgram(*sourceDisas)
		if err != nil {
			panic(err)
		}
		ps := p.Disas()
		if outputDisas == nil {
			for _, inst := range ps {
				fmt.Println(inst)
			}
		} else {
			fd, err := os.Create(*outputDisas)
			if err != nil {
				panic(err)
			}
			defer fd.Close()
			content := strings.Join(ps, "\n")
			fd.WriteString(content)
		}
	}
}
