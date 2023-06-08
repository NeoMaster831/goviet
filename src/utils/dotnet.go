/*
* Parse .NET Utilities / Classes into Go instance.
 */

package utils

import "errors"

/*
public sealed class String : ICloneable, IComparable,
IComparable<string>, IConvertible, IEquatable<string>,
System.Collections.Generic.IEnumerable<char>
*/
type NString struct {
	Str []uint16 // Offset: 0x8
}

func ReadString(hSnap, where uintptr) (NString, error) {
	var ns NString
	var buf uint16 = 0xFFFF
	for ptr := where + 0x8; buf != 0; ptr += 2 {
		ret := RPM(hSnap, ptr, &buf)
		if !ret {
			return ns, errors.New("couldn't read memory")
		}
		ns.Str = append(ns.Str, buf)
	}
	return ns, nil
}
