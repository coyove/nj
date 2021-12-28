package typ

type MarshalType int

const (
	MarshalToString MarshalType = iota
	MarshalToJSON
)

const (
	ForeachContinue = iota
	ForeachBreak
	ForeachDeleteContinue
	ForeachDeleteBreak
)
