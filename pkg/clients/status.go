package clients

type Status int

const (
	Deleted Status = iota
	Active
	Done
)

func (status Status) String() string {
	switch status {
	case Deleted:
		return "deleted"
	case Done:
		return "done"
	default:
		return "active"
	}
}
