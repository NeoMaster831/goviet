package main

import (
	"fmt"
	"pr0j3ct5/goviet/src/utils"
	"pr0j3ct5/goviet/src/winapi"
	"unsafe"
)

const NAME = "gvim.exe"

func main() {
	pid, err := utils.GetPID(NAME)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mBase, err := utils.GetModuleBaseAddress(pid, NAME)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Module Base: 0x%x\n", mBase)

	handle, err := winapi.OpenProcess(0x001FFFFF, 0, pid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var sample uint8
	ret := winapi.ReadProcessMemory(handle, mBase+0x31337, uintptr(unsafe.Pointer(&sample)), uint(unsafe.Sizeof(sample)))

	if ret == 0 {
		fmt.Println("Couldn't read memory")
	} else {
		fmt.Printf("0x%x is the memory address of %x (1 byte)", sample, mBase+0x31337)
	}

}
