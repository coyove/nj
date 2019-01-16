//+build amd64 arm64 ppc64le 386 arm mipsle mips64le

package potatolang

import "unsafe"

// equivalent to: op, opa, opb := op(cruRead64(...))
func cruop(data uintptr, cursor *uint32) (byte, uint16, uint16) {
	addr := uintptr(*cursor) * 4
	*cursor++
	v := *(*uint32)(unsafe.Pointer(data + addr))
	return op(v)
}
