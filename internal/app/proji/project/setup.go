package project

import "database/sql"

// Setup contains necessary informations for the creation of a project.
// Owd is the Origin Working Directory.
type Setup struct {
	Owd          string
	Label        string
	ConfigDir    string
	DatabaseName string
	db           *sql.DB
}

// init initializes the setup struct. Creates a database connection and defines default directores.
func (setup *Setup) init() error {
	// Set original working directory
	if setup.Owd[:len(setup.Owd)-1] != "/" {
		setup.Owd += "/"
	}

	// Connect to database
	connstr := setup.ConfigDir + "/db" + setup.DatabaseName
	db, err := sql.Open("sqlite3", connstr)
	if err != nil {
		return err
	}
	setup.db = db
	return nil
}

// stop cleanly stops the running Setup instance.
// Currently it's only closing its open database connection.
func (setup *Setup) stop() {
	// Close database connection
	if setup.db != nil {
		setup.db.Close()
	}
}

// isLabelSupported checks if the given label is found in the database.
// Returns nil if found, returns error if not found
func (setup *Setup) isLabelSupported() (int, error) {
	stmt, err := setup.db.Prepare("SELECT class_id FROM class_label WHERE label = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	var id int
	err = stmt.QueryRow(setup.Label).Scan(&id)
	return id, err
}
