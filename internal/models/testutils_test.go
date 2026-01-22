package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {

	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")

	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		err2 := db.Close()
		if err2 != nil {
			t.Fatal(err2)
		}
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))

	if err != nil {
		err2 := db.Close()
		if err2 != nil {
			t.Fatal(err2)
		}
		t.Fatal(err)
	}

	// this funcion will be automatically called by go when the current test finishes
	t.Cleanup(func() {
		defer func() {
			err2 := db.Close()
			if err2 != nil {
				t.Fatal(err2)
			}
		}()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))

		if err != nil {
			err2 := db.Close()
			if err2 != nil {
				t.Fatal(err2)
			}
			t.Fatal(err)
		}

	})

	return db
}
