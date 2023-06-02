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

func RPM(hSnap, whereas uintptr, data *interface{}) bool {
	return winapi.ReadProcessMemory(hSnap, whereas, uintptr(unsafe.Pointer(data)), uint(unsafe.Sizeof(data))) != 0
}

func WPM(hSnap, whereas uintptr, data *interface{}) bool {
	return winapi.WriteProcessMemory(hSnap, whereas, uintptr(unsafe.Pointer(data)), uint(unsafe.Sizeof(data))) != 0
}

// Pretty unstable, but I will handle it.
func GetPointerChainValue(hSnap, mBase uintptr, offsets ...int32) (uintptr, error) {

	now := mBase
	var err error = nil

	for i := 0; i < len(offsets)-1; i++ {
		result := winapi.ReadProcessMemory(hSnap, now+uintptr(offsets[i]), uintptr(unsafe.Pointer(&now)), uint(unsafe.Sizeof(now)))
		if result == 0 {
			err = errors.New("couldn't get pointer chain value")
			break
		}
	}

	now += uintptr(offsets[len(offsets)-1])

	return now, err
}
