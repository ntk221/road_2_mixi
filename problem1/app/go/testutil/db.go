package testutil

import (
	"database/sql"
	"fmt"
	"problem1/configs"
	"testing"
)

func OpenDBForTest(t *testing.T) (*sql.DB, error) {
	t.Helper()

	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	t.Cleanup(
		func() { _ = db.Close() },
	)

	return db, nil
}
