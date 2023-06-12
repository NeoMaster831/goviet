/*
* Load kernel API into Go
 */

package winapi

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	MAX_PATH          uint32 = 260
	MAX_MODULE_NAME32 uint32 = 255
)

var (
	kernel32                     = syscall.MustLoadDLL("kernel32.dll")
	winmmDLL                     = syscall.NewLazyDLL("winmm.dll")
	ProcReadProcessMemory        = kernel32.MustFindProc("ReadProcessMemory")
	ProcWriteProcessMemory       = kernel32.MustFindProc("WriteProcessMemory")
	ProcProcess32First           = kernel32.MustFindProc("Process32First")
	ProcProcess32Next            = kernel32.MustFindProc("Process32Next")
	ProcCreateToolhelp32Snapshot = kernel32.MustFindProc("CreateToolhelp32Snapshot")
	ProcOpenProcess              = kernel32.MustFindProc("OpenProcess")
	ProcModule32First            = kernel32.MustFindProc("Module32First")
	ProcModule32Next             = kernel32.MustFindProc("Module32Next")
	ProcVirtualQueryEx           = kernel32.MustFindProc("VirtualQueryEx")
	ProcTimeBeginPeriod          = winmmDLL.NewProc("timeBeginPeriod")
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

type ModuleEntry32 struct {
	DwSize        uint32
	Th32ModuleID  uint32
	Th32ProcessID uint32
	GlblcntUsage  uint32
	ProccntUsage  uint32
	ModBaseAddr   *uint8
	ModBaseSize   uint32
	HModule       uintptr
	SzModule      [MAX_MODULE_NAME32 + 1]uint8
	SzExePath     [MAX_PATH]uint8
}

type Mbi struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	PartitionId       uint16
	RegionSize        uint
	State             uint32
	Protect           uint32
	Type              uint32
}

func CreateToolhelp32Snapshot(dwFlags uint32, th32ProcessID uint32) (ret uintptr) {
	ret, _, _ = ProcCreateToolhelp32Snapshot.Call(
		uintptr(dwFlags),
		uintptr(th32ProcessID),
	)
	return
}

func WriteProcessMemory(hProcess, lpBaseAddress, lpBuffer uintptr, nSize uint) uint8 {
	ret, _, _ := ProcWriteProcessMemory.Call(
		hProcess,
		lpBaseAddress,
		lpBuffer,
		uintptr(nSize),
		0,
	)

	return uint8(ret)
}

func ReadProcessMemory(hProcess, lpBaseAddress, lpBuffer uintptr, nSize uint) uint8 {
	ret, _, _ := ProcReadProcessMemory.Call(
		hProcess,
		lpBaseAddress,
		lpBuffer,
		uintptr(nSize),
		0,
	)

	return uint8(ret)
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

func OpenProcess(dwDesiredAccess uint32, bInheritHandle uint8, dwProcessId uint32) (uintptr, error) {
	ret, _, _ := ProcOpenProcess.Call(
		uintptr(dwDesiredAccess),
		uintptr(bInheritHandle),
		uintptr(dwProcessId),
	)

	if ret == 0 {
		return 0, errors.New("failed to open process")
	}

	return ret, nil
}

func Module32First(hSnap uintptr, entry *ModuleEntry32) bool {
	ret, _, _ := ProcModule32First.Call(
		hSnap,
		uintptr(unsafe.Pointer(entry)),
	)

	return ret != 0
}

func Module32Next(hSnap uintptr, entry *ModuleEntry32) bool {
	ret, _, _ := ProcModule32Next.Call(
		hSnap,
		uintptr(unsafe.Pointer(entry)),
	)

	return ret != 0
}

func VirtualQueryEx(hSnap, lpAddress uintptr, lpBuffer *Mbi, size uint) uint {
	ret, _, _ := ProcVirtualQueryEx.Call(
		hSnap,
		lpAddress,
		uintptr(unsafe.Pointer(lpBuffer)),
		uintptr(size),
	)

	return uint(ret)
}
