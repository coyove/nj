package internal

import (
	"encoding/binary"

	"github.com/coyove/nj/typ"
)

type VByte32 struct {
	Name string
	b    []byte
}

func (p *VByte32) Len() int {
	return len(p.b)
}

func (p *VByte32) Append(idx uint32, line uint32) {
	v := func(v uint64) {
		p.b = append(p.b, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
		n := binary.PutUvarint(p.b[p.Len()-10:], v)
		p.b = p.b[:p.Len()-10+n]
	}
	v(uint64(idx))
	v(uint64(line))
}

func (p *VByte32) Pop() (idx, line uint32) {
	rd := p.b
	a, n := binary.Uvarint(rd)
	b, n2 := binary.Uvarint(rd[n:])
	if n == 0 || n2 == 0 {
		p.b = p.b[:0]
		return
	}
	p.b = p.b[n+n2:]
	return uint32(a), uint32(b)
}

func (p *VByte32) Read(i int) (next int, idx, line uint32) {
	rd := p.b[i:]
	a, n := binary.Uvarint(rd)
	b, n2 := binary.Uvarint(rd[n:])
	if n == 0 || n2 == 0 {
		next = p.Len() + 1
		return
	}
	return i + n + n2, uint32(a), uint32(b)
}

type Packet struct {
	Code []typ.Inst
	Pos  VByte32
}

func (b *Packet) WriteInst(op byte, opa, opb uint16) {
	if opa == opb && op == typ.OpSet {
		return
	}
	b.Code = append(b.Code, typ.Inst{Opcode: op, A: opa, B: int32(opb)})
	if b.Len() >= 4e9 {
		panic("too much code")
	}
}

func (b *Packet) WriteJmpInst(op byte, d int) {
	b.Code = append(b.Code, typ.JmpInst(op, d))
	if b.Len() >= 4e9 {
		panic("too much code")
	}
}

func (b *Packet) WriteLineNum(line uint32) {
	if line == 0 {
		// Debug Code, used to detect a null meta struct
		panic("DEBUG: null line")
	}
	b.Pos.Append(uint32(len(b.Code)), line)
}

func (b *Packet) TruncLast() {
	if len(b.Code) > 0 {
		b.Code = b.Code[:len(b.Code)-1]
	}
}

func (b *Packet) Len() int {
	return len(b.Code)
}

func (b *Packet) LastInst() typ.Inst {
	return b.Code[len(b.Code)-1]
}
