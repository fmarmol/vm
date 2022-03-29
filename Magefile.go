//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fmarmol/basename/pkg/basename"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Test() error {
	return sh.Run("go", "test", "-v", "./pkg/...")
}

// Runs go mod download and then installs the binary.
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "build", "-v", ".")
}

func BuildExamples() error {
	mg.Deps(Build)
	files, err := filepath.Glob("./examples/*.evm")
	if err != nil {
		return err
	}
	for _, file := range files {
		err = sh.Run("./vm", "c", file)
		if err != nil {
			return err
		}
	}
	return nil
}

func Disas() error {
	mg.Deps(BuildExamples)
	files, err := filepath.Glob("./examples/*.vm")
	if err != nil {
		return err
	}
	for _, file := range files {
		fi := basename.ParseFile(file)
		fi.Ext = "disas.evm"
		err = sh.Run("./vm", "disas", file, "-o", fi.FullPath())
		if err != nil {
			return err
		}
	}
	return nil
}

func CleanDisas() error {
	files, err := filepath.Glob("./examples/*.disas.*")
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			return err
		}
		fmt.Println("file:", file, "removed")
	}
	return nil
}

func RunExamples() error {
	mg.Deps(BuildExamples)
	files, err := filepath.Glob("./examples/*.vm")
	if err != nil {
		return err
	}
	for _, file := range files {
		err = sh.Run("./vm", "run", file)
		if err != nil {
			return err
		}
	}
	return nil
}
