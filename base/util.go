package base

type CMap struct {
	Parent *CMap
	M      map[string]int16
}

func NewCMap() *CMap {
	return &CMap{
		M: make(map[string]int16),
	}
}

func (c *CMap) GetRelPosition(key string) int32 {
	m := c
	depth := int32(0)

	for m != nil {
		k, e := m.M[key]
		if e {
			return (depth << 16) | int32(k)
		}

		depth++
		m = m.Parent
	}

	return -1
}
