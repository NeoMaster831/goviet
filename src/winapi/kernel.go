/*
* Load kernel API into Go
 */

package winapi

import "syscall"

var (
	kernel32 = syscall.MustLoadDLL("kernel32.dll")
)
