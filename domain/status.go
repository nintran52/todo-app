package domain

type Status int

const (
	Deleted Status = iota
	Active
	Done
)
