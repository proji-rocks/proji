package item

// Status represents a project status
type Status struct {
	ID        uint   // The status id in storage
	Title     string // The status title
	IsDefault bool   // Is this a default status?
	Comment   string // Short comment describing the status
}

// NewStatus returns a new status
func NewStatus(statusID uint, title, comment string, isDefault bool) *Status {
	return &Status{
		ID:        statusID,
		Title:     title,
		IsDefault: isDefault,
		Comment:   comment,
	}
}
