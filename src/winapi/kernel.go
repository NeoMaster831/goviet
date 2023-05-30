/*
* Load kernel API into Go
 */

package winapi

import "syscall"

var (
	kernel32               = syscall.MustLoadDLL("kernel32.dll")
	ProcReadProcessMemory  = kernel32.MustFindProc("ReadProcessMemory")
	ProcWriteProcessMemory = kernel32.MustFindProc("WriteProcessMemory")
)

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
