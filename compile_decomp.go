package potatolang

import "github.com/coyove/potatolang/parser"

// decompound() will accept a list of atoms, for every compound atom inside,
// it will decompound it into a new temp variable and replace the original one with a Naddr node of this variable.
// for the last compound (if any) decompounded, it will not be saved into a temp variable and used directly (to save some space).
// if we can use r2 register and it has not been touched by the last compound,
// the second last compound (if any) will be saved into r2. In the end we will call OP_RX(..., 2) to transfer.
func (table *symtable) decompound(atoms []*parser.Node, ops []uint16, useR2 bool) (buf packet, err error) {
	// replacedAtoms := []*parser.Node{}
	var lastReplacedAtom, lastlastReplacedAtom struct {
		node, oldnode *parser.Node
		index         int
		lastopPos     int
	}
	buf = newpacket()

	for i, atom := range atoms {
		var yx uint32
		var code packet

		if atom.Type == parser.Ncompound {
			if code, yx, err = table.compileCompoundInto(atom, true, 0); err != nil {
				return
			}
			if table.im != nil {
				atoms[i] = parser.NNode(*table.im)
				table.im = nil
			} else if table.ims != nil {
				atoms[i] = parser.SNode(*table.ims)
				table.ims = nil
			} else {
				// replacedAtoms = append(replacedAtoms, atoms[i])
				lastlastReplacedAtom = lastReplacedAtom
				lastReplacedAtom.oldnode = atom
				atoms[i] = parser.NewNode(parser.Naddr).SetValue(yx)
				buf.Write(code)
				lastReplacedAtom.node = atoms[i]
				lastReplacedAtom.index = i
				lastReplacedAtom.lastopPos = buf.Len() - 1
			}
		}
	}

	if lastReplacedAtom.node != nil {
		_, _, lastReplacedAtom.node.Value = op(buf.data[len(buf.data)-1])
		buf.TruncateLast(1)
		table.sp--
		if ops != nil {
			bop, opa, _ := op(buf.data[len(buf.data)-1])
			idx := uint32(byte(ops[lastReplacedAtom.index]>>8)-OP_R0) / 2
			flag := false

			if flatOpMappingRev[bop] != "" {
				buf.data[len(buf.data)-1] = makeop(bop, idx+1, 0)
				flag = true
			} else if bop == OP_CALL {
				buf.data[len(buf.data)-1] = makeop(OP_CALL, opa, idx+1)
				flag = true
			}

			if flag {
				ops[lastReplacedAtom.index] = OP_NOP
				table.regs[idx].k = false
				table.regs[idx].addr = regA
			}
		}
	}

	if lastlastReplacedAtom.node != nil &&
		ops != nil &&
		!lastlastReplacedAtom.oldnode.WillAffectR2() &&
		!lastReplacedAtom.oldnode.WillAffectR2() && useR2 {
		// r2 trick
		_, _, srcaddr := op(buf.data[lastlastReplacedAtom.lastopPos])
		buf.data[lastlastReplacedAtom.lastopPos] = makeop(OP_R2, srcaddr, 0)
		idx := uint32(byte(ops[lastlastReplacedAtom.index]>>8)-OP_R0) / 2
		if idx != 2 {
			buf.WriteOP(OP_RX, idx, 2)
		}
		ops[lastlastReplacedAtom.index] = OP_NOP
		table.regs[idx].k = false
		table.regs[idx].addr = regA

		flag := false
		lastlastopPos := lastlastReplacedAtom.lastopPos - 1
		lastlastop, lastlasta, _ := op(buf.data[lastlastopPos])

		if flatOpMappingRev[lastlastop] != "" {
			buf.data[lastlastopPos] = makeop(lastlastop, 3, 0)
			flag = true
		} else if lastlastop == OP_CALL {
			buf.data[lastlastopPos] = makeop(OP_CALL, lastlasta, 3)
			flag = true
		}

		if flag {
			// buf.data[lastlastopPos+1] = makeop(OP_NOP, 0, 0)
			buf.data = append(buf.data[:lastlastopPos+1], buf.data[lastlastopPos+2:]...)
		}
	}

	return buf, nil
}
