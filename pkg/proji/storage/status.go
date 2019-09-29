package storage

// Status represents a project status
type Status struct {
	// The status id
	ID uint

	// The status title
	Title string

	// Short comment describing the status.
	Comment string
}

// NewStatus returns a new status
func NewStatus(statusID uint, title, comment string) *Status {
	return &Status{
		ID:      statusID,
		Title:   title,
		Comment: comment,
	}
}
