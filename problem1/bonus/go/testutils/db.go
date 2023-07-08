package testutils

import (
	"database/sql"
	"problem1/configs"
	"testing"
)

func OpenDBForTest(t *testing.T) *sql.DB {
	t.Helper()

	conf := configs.Get()
	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	return db
}
