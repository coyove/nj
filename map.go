package potatolang

import (
	"bytes"
	"io"
)

type Map struct {
	root *node
	size int
}

type node struct {
	left, right *node
	parent      *node
	key, value  Value
}

func (m *Map) leftRotate(x *node) {
	y := x.right
	if y != nil {
		x.right = y.left
		if y.left != nil {
			y.left.parent = x
		}
		y.parent = x.parent
	}

	if x.parent == nil {
		m.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	if y != nil {
		y.left = x
	}
	x.parent = y
}

func (m *Map) rightRotate(x *node) {
	y := x.left
	if y != nil {
		x.left = y.right
		if y.right != nil {
			y.right.parent = x
		}
		y.parent = x.parent
	}
	if x.parent == nil {
		m.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	if y != nil {
		y.right = x
	}
	x.parent = y
}

func (n *node) isRightLeaf() bool { return n.parent.right == n }

func (n *node) isLeftLeaf() bool { return n.parent.left == n }

func (m *Map) splay(x *node) {
	for (x.parent) != nil {
		if x.parent.parent == nil {
			if x.isLeftLeaf() {
				m.rightRotate(x.parent)
			} else {
				m.leftRotate(x.parent)
			}
		} else if x.isLeftLeaf() && x.parent.isLeftLeaf() {
			m.rightRotate(x.parent.parent)
			m.rightRotate(x.parent)
		} else if x.isRightLeaf() && x.parent.isRightLeaf() {
			m.leftRotate(x.parent.parent)
			m.leftRotate(x.parent)
		} else if x.isLeftLeaf() && x.parent.isRightLeaf() {
			m.rightRotate(x.parent)
			m.leftRotate(x.parent)
		} else {
			m.leftRotate(x.parent)
			m.rightRotate(x.parent)
		}
	}
}

func (m *Map) replace(u, v *node) {
	if u.parent == nil {
		m.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	if v != nil {
		v.parent = u.parent
	}
}

func (m *Map) Put(key, value Value) {
	if key.IsNil() {
		return
	}

	if value.IsNil() {
		m.delete(key)
		return
	}

	var z *node = m.root
	var p *node = nil

	for z != nil {
		p = z
		if c := z.key.cmp(key); c == 0 {
			z.value = value
			return
		} else if c < 0 {
			z = z.right
		} else {
			z = z.left
		}
	}

	z = &node{key: key, value: value}
	z.parent = p

	if p == nil {
		m.root = z
	} else if p.key.cmp(z.key) < 0 {
		p.right = z
	} else {
		p.left = z
	}

	m.splay(z)
	m.size++
}

func (m *Map) Get(key Value) Value {
	n := m.find(key)
	if n == nil {
		return Value{}
	}
	m.splay(n)
	return n.value
}

func (m *Map) find(key Value) *node {
	z := m.root
	for z != nil {
		if c := z.key.cmp(key); c < 0 {
			// z.key < key
			z = z.right
		} else if c > 0 {
			// key < z.key
			z = z.left
		} else {
			return z
		}
	}
	return nil
}

func (m *Map) Len() int { return m.size }

func (m *Map) delete(key Value) {
	z := m.find(key)
	if z == nil {
		return
	}

	m.splay(z)

	if z.left == nil {
		m.replace(z, z.right)
	} else if z.right == nil {
		m.replace(z, z.left)
	} else {
		y := subtreeLeftmost(z.right)
		if y.parent != z {
			m.replace(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		m.replace(z, y)
		y.left = z.left
		y.left.parent = y
	}
	m.size--
}

func subtreeLeftmost(u *node) *node {
	for u.left != nil {
		u = u.left
	}
	return u
}

func (n *node) debugString(w io.StringWriter, ident string) {
	if n == nil {
		w.WriteString(ident + "(nil)\n")
		return
	}
	w.WriteString(ident + n.key.String() + "=" + n.value.String() + "\n")
	if n.left == nil && n.right == nil {
		return
	}
	n.right.debugString(w, ident+"  ")
	n.left.debugString(w, ident+"  ")
}

func (m *Map) GoString() string {
	p := bytes.Buffer{}
	m.root.debugString(&p, "")
	return p.String()
}

func Next(m *Map, k Value) (Value, Value) {
	if m.root == nil {
		return Value{}, Value{}
	}

	if k.IsNil() {
		n := subtreeLeftmost(m.root)
		m.splay(n)
		return n.key, n.value
	}

	n := m.find(k)
	if n == nil {
		return Value{}, Value{}
	}
	if n.right == nil { // find upward
	AGAIN:
		next := n.parent
		if next == nil {
			return Value{}, Value{}
		}

		if n.isRightLeaf() {
			n = n.parent
			goto AGAIN
		}

		m.splay(next)
		return next.key, next.value
	}

	next := n.right
	if next.left != nil { // find leftmost
		next = subtreeLeftmost(next.left)
	}

	m.splay(next)
	return next.key, next.value
}
