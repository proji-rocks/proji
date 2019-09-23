package storage

import (
	"github.com/nikoksr/proji/internal/app/proji/class"
)

// Service interface describes the behaviour of a storage service.
type Service interface {
	Close() error
	Save(*class.Class) error
	Load(string) (*class.Class, error)
}
