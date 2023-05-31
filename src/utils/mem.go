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
