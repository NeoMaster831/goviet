/*
* Memory Manager
 */

package utils

import (
	"bytes"
	"errors"
	"pr0j3ct5/goviet/src/winapi"
	"reflect"
	"sync"
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

func String2UInt8(inp string) []uint8 {
	return []byte(inp)
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

func RPM(hSnap, where uintptr, storeat interface{}) bool {
	if reflect.TypeOf(storeat).Kind() != reflect.Ptr {
		return false
	}

	return winapi.ReadProcessMemory(
		hSnap,
		where,
		reflect.ValueOf(storeat).Elem().Addr().Pointer(),
		uint(reflect.TypeOf(storeat).Elem().Size()),
	) != 0
}

func WPM(hSnap, where uintptr, storeat interface{}) bool {
	if reflect.TypeOf(storeat).Kind() != reflect.Ptr {
		return false
	}

	return winapi.WriteProcessMemory(
		hSnap,
		where,
		reflect.ValueOf(storeat).Elem().Addr().Pointer(),
		uint(reflect.TypeOf(storeat).Elem().Size()),
	) != 0
}

// Pretty unstable, but I will handle it.
func Get32BitPointerChainValue(hSnap, where uintptr, offsets ...int32) (uintptr, error) {

	now := where
	var err error = nil

	for i := 0; i < len(offsets)-1; i++ {
		result := winapi.ReadProcessMemory(hSnap, now+uintptr(offsets[i]), uintptr(unsafe.Pointer(&now)), uint(unsafe.Sizeof(now)))
		if result == 0 {
			err = errors.New("couldn't get pointer chain value")
			break
		}
		now &= 0x00000000FFFFFFFF
	}

	now += uintptr(offsets[len(offsets)-1])

	return now, err
}

func Get32BitPtr(hSnap, where uintptr) (uintptr, bool) {
	var ret uintptr
	res := winapi.ReadProcessMemory(
		hSnap,
		where,
		uintptr(unsafe.Pointer(&ret)),
		4,
	)
	return ret, res != 0
}

type AobTask struct {
	A     []byte
	M     string
	Start uintptr
	Size  uint
}

func AsyncScanWorker(hSnap uintptr, wg *sync.WaitGroup, chk *bool, aobj *chan AobTask, storeat *uintptr) {
	for task := range *aobj {
		//fmt.Printf("Running Task @ %x(%x), Found? = %t\n", task.Start, task.Size, *chk)
		if *chk {
			(*wg).Done()
			continue
		}
		patsz := len(task.M)
		for i := uint(0); i < task.Size; i++ {
			found := true
			for j := 0; j < patsz; j++ {
				var res uint8
				_ = RPM(hSnap, task.Start+uintptr(i)+uintptr(j), &res)
				if task.M[j] != '?' && task.A[j] != res {
					found = false
					break
				}
			}
			if found {
				*chk = true
				*storeat = task.Start + uintptr(i)
			}
		}
		(*wg).Done()
	}
}

func AsyncScanRegion(hSnap uintptr, pattern []uint8, mask string, begin uintptr, size uint, Chan *chan AobTask, wg *sync.WaitGroup) {
	var mbi winapi.Mbi = winapi.Mbi{RegionSize: 0}
	for cur := begin; cur < begin+uintptr(size); cur += uintptr(mbi.RegionSize) {
		if winapi.VirtualQueryEx(hSnap, cur, &mbi, uint(unsafe.Sizeof(mbi))) == 0 ||
			mbi.State != 0x1000 ||
			mbi.Protect&0x100 != 0 ||
			mbi.Protect&0x01 != 0 {
			continue
		}

		wg.Add(1)
		*Chan <- AobTask{pattern[:], mask, cur, mbi.RegionSize}
	}
}
