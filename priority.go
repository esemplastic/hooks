package hooks

type Priority int

const (
	Idle Priority = iota * 100
	BelowNormal
	Normal
	AboveNormal
	High
	Realtime
)
