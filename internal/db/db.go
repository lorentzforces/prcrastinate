package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"prcrastinate/internal/platform"

	_ "github.com/mattn/go-sqlite3"
)

const dbVersion = 1

const defaultDbPath = "~/.local/share/prcrastinate/data.db"

func open(givenPath string) (*sql.DB, error) {
	path, _, _ := platform.GetDefaultablePath(givenPath, defaultDbPath)

	if err := os.MkdirAll(filepath.Dir(path), 0o_777); err != nil {
		platform.FailOut(fmt.Sprintf(
			"Could not create directory structure for db file: %s\n%s",
			path,
			err.Error(),
		))
	}

	connectionString := "file:" + path + "?_foreign_keys=on"
	database, err := sql.Open("sqlite3", connectionString)
	return database, err
}

func createSchema(db *sql.DB) error {
	_, err := db.Exec(tableVersion)
	return err
}

func verifyVersion(db *sql.DB) error {
	row := db.QueryRow("SELECT number FROM version ORDER BY number DESC LIMIT 1")
	var actualVersion int
	if err := row.Scan(&actualVersion); err != nil {
		return err
	}

	if actualVersion != dbVersion {
		return errDbVersionMismatch
	}
	return nil
}

var errDbVersionMismatch = fmt.Errorf(
	"Stored data version does not match the app. Pull fresh data to refresh the data store.")

//go:embed scripts/table_version.sql
var tableVersion string
