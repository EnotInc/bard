package calls

type Call int

const (
	_ Call = iota
	None
	PurgeCache
	OpenFile
)
