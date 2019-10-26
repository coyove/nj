package potatolang

import (
	"github.com/coyove/potatolang/parser"
)

func (table *symtable) decompound(atoms []*parser.Node) (buf packet, err error) {
	return table._decompound(atoms, true)
}

func (table *symtable) decompoundWithoutA(atoms []*parser.Node) (buf packet, err error) {
	return table._decompound(atoms, false)
}

// decompound() will accept a list of atoms, for every compound atom inside,
// it will decompound it into a new temp variable and replace the original one with a Naddr node of this variable.
// for the last compound (if any) decompounded, it will not be saved into a temp variable
// and used directly (to save some space) if useA == true.
func (table *symtable) _decompound(atoms []*parser.Node, useA bool) (buf packet, err error) {
	buf = newpacket()

	var lastCompound struct {
		n *parser.Node
		i int
	}

	for i, atom := range atoms {
		if atom == nil {
			break
		}

		var yx uint16
		var code packet

		if atom.Type == parser.Ncompound {
			if code, yx, err = table.compileCompoundInto(atom, true, 0); err != nil {
				return
			}

			atoms[i] = parser.NewNode(parser.Naddr).SetValue(yx)
			buf.Write(code)

			lastCompound.n = atom
			lastCompound.i = i
		}
	}

	if lastCompound.n != nil && useA {
		_, old, opb := op(buf.data[len(buf.data)-1])
		buf.TruncateLast(1)
		table.returnAddress(old)
		atoms[lastCompound.i].SetValue(opb)

		//if len(buf.data) > 0 {
		//	opcode, opa, opk := op(buf.data[len(buf.data)-1])
		//	if opcode == OpSet && opa == regA && opb == regA {
		//		buf.TruncateLast(1)
		//		atoms[lastCompound.i].SetValue(opk)
		//	}
		//}
	}

	return
}
