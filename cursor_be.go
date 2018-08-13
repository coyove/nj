//+build mips64 ppc64 s390x

package potatolang

import "unsafe"

// equivalent to: op, opa, opb := op(cruRead64(...)) in big endian mode
func cruop(data uintptr, cursor *uint32) (byte, uint32, uint32) {
	addr := uintptr(*cursor) * 8
	*cursor++
	return *(*byte)(unsafe.Pointer(data + addr)),
		*(*uint32)(unsafe.Pointer(data + addr)) & 0x00ffffff,
		*(*uint32)(unsafe.Pointer(data + addr + 4))
}
