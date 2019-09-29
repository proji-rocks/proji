package storage

// Service interface describes the behaviour of a storage service.
type Service interface {
	// Close the interface.
	Close() error

	// SaveClass saves a class to storage.
	SaveClass(class *Class) error

	// LoadClass loads a class from storage by its ID.
	LoadClass(classID uint) (*Class, error)

	// LoadAllClasses loads all available classes from storage.
	LoadAllClasses() ([]*Class, error)

	// RemoveClass removes a class from storage.
	RemoveClass(classID uint) error

	// SaveProject saves a project to storage.
	SaveProject(proj *Project) error

	// LoadProject loads a project from storage by its ID.
	LoadProject(projectID uint) (*Project, error)

	// LoadAllProjects returns a list of all projects in storage.
	LoadAllProjects() ([]*Project, error)

	// LoadProjectID loads the ID of a project.
	LoadProjectID(installPath string) (uint, error)

	// LoadClassIDByLabel loads the ID of a class by its label.
	LoadClassIDByLabel(label string) (uint, error)

	// UpdateProjectStatus updates the status of a given project in storage.
	UpdateProjectStatus(projectID, statusID uint) error

	// UpdateProjectLocation updates the location of a project in storage.
	UpdateProjectLocation(projectID uint, installPath string) error

	// UntrackProject removes a project from storage.
	RemoveProject(projectID uint) error

	// SaveStatus adds a new status to storage.
	SaveStatus(status *Status) error

	// UpdateStatus updates a status in storage.
	UpdateStatus(statusID uint, title, comment string) error

	// LoadStatus loads a status from storage by its ID.
	LoadStatus(statusID uint) (*Status, error)

	// LoadAllStatuses returns a list of all statuses in storage.
	LoadAllStatuses() ([]*Status, error)

	// LoadStatusID loads the ID of a given status from storage.
	LoadStatusID(title string) (uint, error)

	// RemoveStatus removes an existing status from storage.
	RemoveStatus(statusID uint) error
}
