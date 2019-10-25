package sqlite_test

import (
	"os"
	"testing"

	"github.com/nikoksr/proji/pkg/proji/storage/item"

	"github.com/nikoksr/proji/pkg/proji/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestClass(t *testing.T) {
	const numOfDefaultClasses = 1
	dbPath := "/tmp/proji.sqlite3"
	svc, err := sqlite.New(dbPath)
	defer os.Remove(dbPath)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	className := "testclass1"
	classLabel := "tc1"
	class := item.NewClass(className, classLabel)
	assert.NotNil(t, class)

	// Test SaveClass
	err = svc.SaveClass(class)
	assert.NoError(t, err)

	// Test success of SaveClass by loading the class' ID
	classID, err := svc.LoadClassIDByLabel(classLabel)
	assert.NoError(t, err)
	assert.NotEqual(t, classID, 0)
	class.ID = classID

	// Test success of LoadClassIDByLabel by loading the class again by its ID
	testClass, err := svc.LoadClass(classID)
	assert.NoError(t, err)
	assert.Equal(t, class, testClass)

	// Save more classes and try to load them all
	goodClasses := make([]*item.Class, 0)
	goodClasses = append(goodClasses, item.NewClass("testclass2", "tc2"))
	goodClasses = append(goodClasses, item.NewClass("testclass3", "tc3"))
	goodClasses = append(goodClasses, item.NewClass("testclass4", "tc4"))

	for _, goodClass := range goodClasses {
		err = svc.SaveClass(goodClass)
		assert.NoError(t, err)
		id, err := svc.LoadClassIDByLabel(goodClass.Label)
		assert.NoError(t, err)
		assert.NotEqual(t, id, 0)
		goodClass.ID = id
	}

	goodClasses = append([]*item.Class{testClass}, goodClasses...)
	classes, err := svc.LoadAllClasses()
	assert.NoError(t, err)
	assert.NotNil(t, classes)
	assert.Equal(t, len(goodClasses), len(classes)-numOfDefaultClasses) // Subtract the number of default classes

	for idx, class := range classes {
		assert.Equal(t, goodClasses[idx], class)
	}

	// Test unique constraint for name and label
	badClasses := []*item.Class{
		&item.Class{
			Name:  className,
			Label: classLabel,
		},
		&item.Class{
			Name:  className,
			Label: "x",
		},
		&item.Class{
			Name:  "myNewClass",
			Label: classLabel,
		},
	}

	for _, badClass := range badClasses {
		err = svc.SaveClass(badClass)
		assert.Error(t, err)
	}

	// Try to remove all classes
	for _, class := range classes {
		err = svc.RemoveClass(class.ID)
		assert.NoError(t, err)
	}
}

func TestProject(t *testing.T) {
	dbPath := "/tmp/proji.sqlite3"
	svc, err := sqlite.New(dbPath)
	defer os.Remove(dbPath)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	projectsInMem := make([]*item.Project, 0)

	// Create a temporary class
	class := item.NewClass("testclass", "tc")
	err = svc.SaveClass(class)
	assert.NoError(t, err)

	// Load the default status 'active'
	status, err := svc.LoadStatus(1)
	assert.NoError(t, err)

	projectName := "test-proj1"
	basePath := "/tmp/"
	projectPath := basePath + projectName
	proj := item.NewProject(0, projectName, projectPath+projectName, class, status)

	// Test SaveProject
	err = svc.SaveProject(proj)
	assert.NoError(t, err)

	// Should fail because it already exists
	err = svc.SaveProject(proj)
	assert.Error(t, err)

	// Test the project ID
	id, err := svc.LoadProjectID(proj.InstallPath)
	assert.NoError(t, err)
	proj.ID = id

	// Try to load the project and compare it to the in memory one
	loadedProj, err := svc.LoadProject(id)
	assert.NoError(t, err)
	assert.Equal(t, loadedProj.ID, id)
	assert.Equal(t, proj, loadedProj)
	proj = nil

	// Update the status and location of the project
	newStatusID := uint(1)
	err = svc.UpdateProjectStatus(loadedProj.ID, newStatusID)
	assert.NoError(t, err)

	newProjPath := "/test"
	err = svc.UpdateProjectLocation(loadedProj.ID, newProjPath)
	assert.NoError(t, err)

	loadedProj, err = svc.LoadProject(loadedProj.ID)
	assert.NoError(t, err)
	assert.Equal(t, newStatusID, loadedProj.Status.ID)
	assert.Equal(t, newProjPath, loadedProj.InstallPath)
	projectsInMem = append(projectsInMem, loadedProj)

	// Add another project and try to load both of them
	proj2 := item.NewProject(0, "test-proj2", basePath+"test-proj2", class, status)
	err = svc.SaveProject(proj2)
	assert.NoError(t, err)

	id, err = svc.LoadProjectID(proj2.InstallPath)
	assert.NoError(t, err)
	proj2.ID = id

	projectsInMem = append(projectsInMem, proj2)

	projects, err := svc.LoadAllProjects()
	assert.NoError(t, err)
	assert.Equal(t, len(projectsInMem), len(projects))

	for idx, proj := range projects {
		assert.Equal(t, projectsInMem[idx], proj)

		// Try to remove the project
		err = svc.RemoveProject(proj.ID)
		assert.NoError(t, err)
	}
}

func TestStatus(t *testing.T) {
	dbPath := "/tmp/proji.sqlite3"
	svc, err := sqlite.New(dbPath)
	defer os.Remove(dbPath)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	status := item.NewStatus(0, "test", "This is a test status.")
	badStatus := item.NewStatus(1, "active", "This status should already exist.")

	// Try to save the new status; should be successful
	err = svc.SaveStatus(status)
	assert.NoError(t, err)

	// Try to save a status that already exists; should fail
	err = svc.SaveStatus(badStatus)
	assert.Error(t, err)
	badStatus = nil

	// Load the status ID
	id, err := svc.LoadStatusID(status.Title)
	assert.NoError(t, err)
	status.ID = id

	// Load and compare the status
	loadedStatus, err := svc.LoadStatus(status.ID)
	assert.NoError(t, err)
	assert.Equal(t, status, loadedStatus)

	// Update a status
	newTitle := "updated-test"
	newComment := "This is an updated test."
	err = svc.UpdateStatus(status.ID, newTitle, newComment)
	assert.NoError(t, err)

	// Reload and compare again
	status, err = svc.LoadStatus(status.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, loadedStatus, status)
	assert.Equal(t, newTitle, status.Title)
	assert.Equal(t, newComment, status.Comment)
	status.ID = id

	// Add another status, then try to load and remove all
	status2 := item.NewStatus(0, "test2", "This is the second test status.")
	err = svc.SaveStatus(status2)
	assert.NoError(t, err)
	id, err = svc.LoadStatusID(status2.Title)
	assert.NoError(t, err)
	status2.ID = id

	statusesInMem := []*item.Status{status, status2}
	statuses, err := svc.LoadAllStatuses()
	assert.NoError(t, err)

	// Exclude the 5 default statuses
	for idx, s := range statuses[5:] {
		assert.Equal(t, statusesInMem[idx], s)

		// Try to remove the status
		err = svc.RemoveStatus(s.ID)
		assert.NoError(t, err)
	}
}
