package storage

// Service interface describes the behaviour of a storage service.
type Service interface {
	// Close the interface.
	Close() error

	// Save a class to a storage.
	SaveClass(class *Class) error

	// Load a class from storage by its name.
	LoadClassByName(name string) (*Class, error)

	// Load a class from storage by its ID.
	LoadClassByID(id uint) (*Class, error)

	// Load the ID for a given class
	LoadClassID(name string) (uint, error)

	// Load all classes from storage
	LoadAllClasses() ([]*Class, error)

	// Remove a class from storage.
	RemoveClass(name string) error

	// Checks if a given label exists in the storage. Returns the corresponding ID if true and an error if not.
	DoesLabelExist(label string) (uint, error)

	// TrackProject adds a project to the database.
	TrackProject(proj *Project) error
}
