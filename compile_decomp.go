package potatolang

import (
	"github.com/coyove/potatolang/parser"
)

// decompound() will accept a list of atoms, for every compound atom inside,
// it will decompound it into a new temp variable and replace the original one with a Naddr node of this variable.
// for the last compound (if any) decompounded, it will not be saved into a temp variable and used directly (to save some space).
// if we can use r2 register and it has not been touched by the last compound,
// the second last compound (if any) will be saved into r2. In the end we will call OP_RX(..., 2) to transfer.
func (table *symtable) decompound(atoms []*parser.Node, dontUseA ...bool) (buf packet, err error) {
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
			if code, yx, err = table.compileCompoundInto(atom, true, 0, false); err != nil {
				return
			}

			atoms[i] = parser.NewNode(parser.Naddr).SetValue(yx)
			buf.Write(code)

			lastCompound.n = atom
			lastCompound.i = i
		}
	}

	if lastCompound.n != nil && len(dontUseA) == 0 {
		_, _, opb := op(buf.data[len(buf.data)-1])
		buf.TruncateLast(1)
		table.vp--
		atoms[lastCompound.i].SetValue(opb)
	}

	return
}
