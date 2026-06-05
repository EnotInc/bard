package calls

type Call int

const (
	_ Call = iota
	None
	Rezise // TODO: implement
	PurgeCache
)
