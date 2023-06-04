package main

import (
	"fmt"
	"pr0j3ct5/goviet/src/utils"
	"pr0j3ct5/goviet/src/winapi"
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

	var sample uint32
	const nullp = uintptr(0x30E800)
	ret := utils.RPM(handle, mBase+nullp, &sample)
	//ret := winapi.ReadProcessMemory(handle, mBase+0x31337, uintptr(unsafe.Pointer(&sample)), uint(unsafe.Sizeof(sample)))

	if !ret {
		fmt.Println("Couldn't read memory")
		return
	}

	type MyStruct struct {
		value [100]byte
	}
	var towrite MyStruct
	copy(towrite.value[:], utils.String2UInt8("OG for a fee, stay sippin'"))

	//fmt.Println("Size of towrite:", reflect.TypeOf(towrite).Size())
	fmt.Printf("0x%x is the original value, address of %x\n", sample, mBase+nullp)
	ret = utils.WPM(handle, mBase+nullp, &towrite)
	if !ret {
		fmt.Println("Couldn't write memory")
		return
	}

	ret = utils.RPM(handle, mBase+nullp, &sample)
	if !ret {
		fmt.Println("Some shit happened lol")
		return
	}
	fmt.Printf("0x%x is the new value, address of %x\n", sample, mBase+nullp)

}
