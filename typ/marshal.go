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

func (m MarshalType) String() string {
	switch m {
	case MarshalToString:
		return "str"
	case MarshalToStringNonRec:
		return "strnorec"
	case MarshalToJSON:
		return "json"
	}
	return "unknown"
}
