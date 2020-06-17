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
	UpdateProjectLocation(projectID uint, installPath string) error // UpdateProjectLocation updates the location of a project in storage.
	RemoveProject(projectID uint) error                             // RemoveProject removes a project from storage.
}
