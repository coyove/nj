//+build amd64 arm64 ppc64le 386 arm mipsle mips64le

package potatolang

import "unsafe"

// equivalent to: op, opa, opb := op(cruRead64(...))
func cruop(data uintptr, cursor *uint32) (byte, uint32, uint32) {
	addr := uintptr(*cursor) * 8
	*cursor++
	return *(*byte)(unsafe.Pointer(data + addr + 7)),
		*(*uint32)(unsafe.Pointer(data + addr + 4)) & 0x00ffffff,
		*(*uint32)(unsafe.Pointer(data + addr))
}
