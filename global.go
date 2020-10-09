package potatolang

import "encoding/binary"

type DataStack struct {
	Data *[]byte
}

func (g DataStack) Alloc(n int) (ptr uint64, data []byte) {
	alignedLength := n
	if alignedLength/8*8 != alignedLength {
		alignedLength += 8 - alignedLength%8
	}

	tmp := [10]byte{}
	*g.Data = append(*g.Data, tmp[:binary.PutUvarint(tmp[:], uint64(n))]...)

	start := len(*g.Data)
	if alignedLength > len(*g.Data) {
		for i := 0; i < alignedLength; i += 8 {
			*g.Data = append(*g.Data, 0, 0, 0, 0, 0, 0, 0, 0)
		}
	} else {
		*g.Data = append(*g.Data, (*g.Data)[:alignedLength]...)
	}

	return 0x7<<48 | uint64(start), (*g.Data)[start : start+n]
}
