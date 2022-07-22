package internal

import (
	"bytes"
	"encoding/binary"
	"unicode/utf8"

	"github.com/coyove/nj/typ"
)

const (
	MaxStackSize   = 0x0fffffff
	FlagTailCall   = 0x80000000
	FlagNativeCall = 0x40000000
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

func (b *Packet) check() {
	if b.Len() >= MaxStackSize-1 {
		panic("too much code")
	}
}

func (b *Packet) WriteInst(op byte, opa, opb uint16) {
	if op == typ.OpSet {
		if opa == opb {
			/*
				form:
					set v v
			*/
			return
		}
		if opb == typ.RegA && len(b.Code) > 0 {
			/*
				form:
					load subject key $a
					set dest $a
				into:
					load subject key dest
			*/
			x := &b.Code[len(b.Code)-1]
			if x.Opcode == typ.OpLoad && x.C == typ.RegA {
				x.C = opa
				return
			}
		}
		if opb == typ.RegA && len(b.Code) > 0 {
			/*
				form:
					add v num
					set v $a
				into:
					inc v num
				note that 'add num v' is not optimizable because 'add' also applies on strings
			*/
			x := &b.Code[len(b.Code)-1]
			if x.Opcode == typ.OpAdd && x.A == opa {
				x.Opcode = typ.OpInc
				return
			}
		}
	}
	b.Code = append(b.Code, typ.Inst{Opcode: op, A: opa, B: opb})
	b.check()
}

func (b *Packet) WriteInst3(op byte, opa, opb, opc uint16) {
	b.Code = append(b.Code, typ.Inst{Opcode: op, A: opa, B: opb, C: opc})
	b.check()
}

func (b *Packet) WriteJmpInst(op byte, d int) {
	b.Code = append(b.Code, typ.JmpInst(op, d))
	b.check()
}

func (b *Packet) WriteLineNum(line uint32) {
	if line == 0 {
		ShouldNotHappen()
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

func (b *Packet) Copy() *Packet {
	b2 := *b
	b2.Code = append([]typ.Inst{}, b.Code...)
	return &b2
}

type LimitedBuffer struct {
	Limit int
	bytes.Buffer
}

func (w *LimitedBuffer) Write(b []byte) (int, error) {
	if w.Limit > 0 {
		if w.Len()+len(b) > w.Limit {
			if _, err := w.Buffer.Write(b[:w.Limit-w.Len()]); err != nil {
				return 0, err
			}
			return len(b), nil
		}
	}
	return w.Buffer.Write(b)
}

func (w *LimitedBuffer) WriteString(b string) (int, error) {
	if w.Limit > 0 {
		if w.Len()+len(b) > w.Limit {
			if _, err := w.Buffer.WriteString(b[:w.Limit-w.Len()]); err != nil {
				return 0, err
			}
			return len(b), nil
		}
	}
	return w.Buffer.WriteString(b)
}

func (w *LimitedBuffer) WriteByte(b byte) error {
	if w.Limit > 0 {
		if w.Len()+1 > w.Limit {
			return nil
		}
	}
	return w.Buffer.WriteByte(b)
}

func (w *LimitedBuffer) WriteRune(b rune) (int, error) {
	if w.Limit > 0 {
		sz := utf8.RuneLen(b)
		if w.Len()+sz > w.Limit {
			return sz, nil
		}
	}
	return w.Buffer.WriteRune(b)
}
