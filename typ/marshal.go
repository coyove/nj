package typ

type MarshalType int

const (
	MarshalToString MarshalType = iota
	MarshalToStringNonRec
	MarshalToJSON
)

func (m MarshalType) NoRec() MarshalType {
	if m == MarshalToJSON {
		return m
	}
	return MarshalToStringNonRec
}
