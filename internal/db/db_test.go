package db

import (
	"database/sql"
	"errors"
	"testing"
)

func TestDbCreation(t *testing.T) {
	conn := createTestDb(t)
	if err := createSchema(conn); err != nil {
		t.Fatal(err)
	}
}

func TestVersionVerify(t *testing.T) {
	conn := createTestDb(t)
	if err := createSchema(conn); err != nil {
		t.Fatal(err)
	}
	// -1 should never match
	insertVersion := "INSERT INTO version (number) VALUES (-1)"
	if _, err := conn.Exec(insertVersion); err != nil {
		t.Fatal(err)
	}

	err := verifyVersion(conn)
	if err == nil {
		t.Fatalf(
			"Database version (-1) should provoke mismatch with expected version, but did not.")
	}
	if !errors.Is(err, errDbVersionMismatch) {
		t.Fatalf("Expected errDbVersionMismatch error, but was given an error instead: %s", err)
	}
}

func createTestDb(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return db
}
