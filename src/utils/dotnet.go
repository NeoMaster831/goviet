/*
* Parse .NET Utilities / Classes into Go instance.
 */

package utils

type NString struct {
	header byte
}

func ReadString(hSnap, where uintptr) (NString, error) {
	return NString{0x0}, nil
}
