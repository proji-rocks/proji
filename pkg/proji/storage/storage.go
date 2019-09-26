package storage

// Service interface describes the behaviour of a storage service.
type Service interface {
	// Close the interface.
	Close() error

	// SaveClass saves a class to storage.
	SaveClass(class *Class) error

	// LoadClassByName loads a class from storage by its name.
	LoadClassByName(name string) (*Class, error)

	// LoadClassByID loads a class from storage by its ID.
	LoadClassByID(id uint) (*Class, error)

	// LoadClassID loads the ID of a given class from storage.
	LoadClassID(name string) (uint, error)

	// LoadAllClasses loads all available classes from storage.
	LoadAllClasses() ([]*Class, error)

	// RemoveClass removes a class from storage.
	RemoveClass(name string) error

	// DoesLabelExist checks if a given label exists in storage. Returns the corresponding ID if true and an error if not.
	DoesLabelExist(label string) (uint, error)

	// TrackProject adds a project to storage.
	TrackProject(proj *Project) error

	// UntrackProject removes a project from storage.
	UntrackProject(id uint) error

	// UpdateProjectStatus updates the status of a given project in storage.
	UpdateProjectStatus(projectID, statusID uint) error

	// UpdateProjectLocation updates the location of a project in storage.
	UpdateProjectLocation(projectID uint, installPath string) error

	// LoadProjectID loads the ID of a given project from storage.
	LoadProjectID(path string) (uint, error)

	// ListProjects returns a list of all projects in storage.
	ListProjects() ([]*Project, error)

	// AddStatus adds a new status to storage.
	AddStatus(status *Status) error

	// UpdateStatus updates a status in storage.
	UpdateStatus(status *Status) error

	// RemoveStatus removes an existing status from storage.
	RemoveStatus(statusID uint) error

	// LoadStatusByTitle loads a status from storage by its title.
	LoadStatusByTitle(title string) (*Status, error)

	// LoadStatusByID loads a status from storage by its ID.
	LoadStatusByID(id uint) (*Status, error)

	// LoadStatusID loads the ID of a given status from storage.
	LoadStatusID(title string) (uint, error)

	// ListAvailableProjectStatuses returns a list of all statuses in storage.
	ListAvailableStatuses() ([]*Status, error)
}
