package class

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"
	"github.com/spf13/viper"
)

// ListAll lists all classes available in the database
func ListAll() error {
	DBDir := helper.GetConfigDir() + "/db/"
	databaseName, ok := viper.Get("database.name").(string)

	if !ok {
		return errors.New("could not read database name from config file")
	}

	db, err := sql.Open("sqlite3", DBDir+databaseName)
	if err != nil {
		return err
	}
	defer db.Close()

	names, err := db.Query("SELECT name FROM class ORDER BY name")
	if err != nil {
		return err
	}
	defer names.Close()

	for names.Next() {
		var name string
		names.Scan(&name)
		fmt.Println(name)
	}

	return nil
}
