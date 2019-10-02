package storage

// Status represents a project status
type Status struct {
	ID      uint   // The status id in storage
	Title   string // The status title
	Comment string // Short comment describing the status
}

// NewStatus returns a new status
func NewStatus(statusID uint, title, comment string) *Status {
	return &Status{
		ID:      statusID,
		Title:   title,
		Comment: comment,
	}
}
