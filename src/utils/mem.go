/*
* Memory Manager
 */

package utils

import (
	"bytes"
	"errors"
	"pr0j3ct5/goviet/src/winapi"
	"unsafe"
)

func GetPID(pName string) (uint32, error) {

	// TH32_SNAPALL = 15
	handle := winapi.CreateToolhelp32Snapshot(15, 0)
	var entry winapi.ProcessEntry32
	entry.DwSize = uint32(unsafe.Sizeof(entry))

	for valid := winapi.Process32First(handle, &entry); valid; valid = winapi.Process32Next(handle, &entry) {
		if UInt82String(entry.SzExeFile[:]) == pName {
			return entry.Th32ProcessID, nil
		}
	}

	return 0, errors.New("process not found")
}

func UInt82String(inp []uint8) string {
	return string(inp[:bytes.Index(inp, []uint8{0})])
}

func GetModuleBaseAddress(pid uint32, mName string) (uintptr, error) {

	// TH32CS_SNAPMODULE | TH32CS_SNAPMODULE32 = 0x10 | 0x8
	handle := winapi.CreateToolhelp32Snapshot(0x18, pid)
	var entry winapi.ModuleEntry32
	entry.DwSize = uint32(unsafe.Sizeof(entry))

	for valid := winapi.Module32First(handle, &entry); valid; valid = winapi.Module32Next(handle, &entry) {
		if UInt82String(entry.SzModule[:]) == mName {
			return uintptr(unsafe.Pointer(entry.ModBaseAddr)), nil
		}
	}

	return 0, errors.New("module not found")
}
