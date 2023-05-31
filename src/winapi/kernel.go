/*
* Load kernel API into Go
 */

package winapi

import (
	"syscall"
	"unsafe"
)

const (
	MAX_PATH uint32 = 260
)

var (
	kernel32                     = syscall.MustLoadDLL("kernel32.dll")
	ProcReadProcessMemory        = kernel32.MustFindProc("ReadProcessMemory")
	ProcWriteProcessMemory       = kernel32.MustFindProc("WriteProcessMemory")
	ProcProcess32First           = kernel32.MustFindProc("Process32First")
	ProcProcess32Next            = kernel32.MustFindProc("Process32Next")
	ProcCreateToolhelp32Snapshot = kernel32.MustFindProc("CreateToolhelp32Snapshot")
)

type ProcessEntry32 struct {
	DwSize              uint32
	CntUsage            uint32
	Th32ProcessID       uint32
	Th32DefaultHeapID   uintptr
	Th32ModuleID        uint32
	CntThreads          uint32
	Th32ParentProcessID uint32
	PcPriClassBase      int32
	DwFlags             uint32
	SzExeFile           [MAX_PATH]uint8
}

func CreateToolhelp32Snapshot(dwFlags uint32, th32ProcessID uint32) (ret uintptr) {
	ret, _, _ = ProcCreateToolhelp32Snapshot.Call(
		uintptr(dwFlags),
		uintptr(th32ProcessID),
	)
	return
}

func WriteProcessMemory(hProcess, lpBaseAddress, lpBuffer uintptr, nSize uint) (uintptr, error) {
	ret, _, err := ProcWriteProcessMemory.Call(
		hProcess,
		lpBaseAddress,
		lpBuffer,
		uintptr(nSize),
		0,
	)

	if err.Error() != "The operation completed successfully." {
		return 0, err
	}

	return ret, nil
}

func ReadProcessMemory(hProcess, lpBaseAddress, lpBuffer uintptr, nSize uint) (uintptr, error) {
	ret, _, err := ProcReadProcessMemory.Call(
		hProcess,
		lpBaseAddress,
		lpBuffer,
		uintptr(nSize),
		0,
	)

	if err.Error() != "The operation completed successfully." {
		return 0, err
	}

	return ret, nil
}

func Process32First(hSnap uintptr, entry *ProcessEntry32) bool {

	ret, _, _ := ProcProcess32First.Call(
		hSnap,
		uintptr(unsafe.Pointer(entry)),
	)

	return ret != 0
}

func Process32Next(hSnap uintptr, entry *ProcessEntry32) bool {

	ret, _, _ := ProcProcess32Next.Call(
		hSnap,
		uintptr(unsafe.Pointer(entry)),
	)

	return ret != 0
}
