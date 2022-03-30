package main

import "syscall"

func AllocMem(size uint32) ([]byte, error) {
	return syscall.Mmap(-1, 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
}

func FreeMem(data []byte) error {
	return syscall.Munmap(data)
}
