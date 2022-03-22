package main

// DO NOT USE HASH MAP FOR NOW, KEEP THIS CODE FOR LATER
type MemRegister map[uint32]uint32

// func (m *MemRegister) Register(ptr uint32, size uint32) (uint32, Err) {
// 	_, ok := (*m)[ptr]
// 	if ok {
// 		return NULL, Err_AddrAlreadyAllocated
// 	}
// 	return ptr, OK
// }

// func (m *MemRegister) Unregister(ptr uint32) {
// 	delete(*m, ptr)
// }
