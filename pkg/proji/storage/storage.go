package storage

import "github.com/nikoksr/proji/pkg/proji/storage/item"

// Service interface describes the behaviour of a storage service.
type Service interface {
	Close() error                                                   // Close the interface.
	SaveClass(class *item.Class) error                              // SaveClass saves a class to storage.
	LoadClass(classID uint) (*item.Class, error)                    // LoadClass loads a class from storage by its ID.
	LoadAllClasses() ([]*item.Class, error)                         // LoadAllClasses loads all available classes from storage.
	RemoveClass(classID uint) error                                 // RemoveClass removes a class from storage.
	SaveProject(proj *item.Project) error                           // SaveProject saves a project to storage.
	LoadProject(projectID uint) (*item.Project, error)              // LoadProject loads a project from storage by its ID.
	LoadAllProjects() ([]*item.Project, error)                      // LoadAllProjects returns a list of all projects in storage.
	LoadProjectID(installPath string) (uint, error)                 // LoadProjectID loads the ID of a project.
	LoadClassIDByLabel(label string) (uint, error)                  // LoadClassIDByLabel loads the ID of a class by its label.
	UpdateProjectStatus(projectID, statusID uint) error             // UpdateProjectStatus updates the status of a given project in storage.
	UpdateProjectLocation(projectID uint, installPath string) error // UpdateProjectLocation updates the location of a project in storage.
	RemoveProject(projectID uint) error                             // RemoveProject removes a project from storage.
	SaveStatus(status *item.Status) error                           // SaveStatus adds a new status to storage.
	UpdateStatus(statusID uint, title, comment string) error        // UpdateStatus updates a status in storage.
	LoadStatus(statusID uint) (*item.Status, error)                 // LoadStatus loads a status from storage by its ID.
	LoadAllStatuses() ([]*item.Status, error)                       // LoadAllStatuses returns a list of all statuses in storage.
	LoadStatusID(title string) (uint, error)                        // LoadStatusID loads the ID of a given status from storage.
	RemoveStatus(statusID uint) error                               // RemoveStatus removes an existing status from storage.
}
