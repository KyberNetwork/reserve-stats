package migratedb

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
)

const (
	postgresHost     = "127.0.0.1"
	postgresPort     = 5432
	postgresUser     = "reserve_stats"
	postgresPassword = "reserve_stats"
	postgresDatabase = "reserve_stats"
	dbPath           = "testdata/dev_users.db"
)

func newTestMigrateDB() (*DBMigration, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
	)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return NewMigrateStorage(dbPath, db)
}

func tearDown(t *testing.T, dbMigration *DBMigration) {
	if err := dbMigration.DeleteAllTables(); err != nil {
		t.Fatal(err)
	}
	// assert.Nil(t, dbMigration.DeleteAllTables(), "test db should be tear down succesfully.")
}

func TestMigrateDB(t *testing.T) {
	dbMigration, err := newTestMigrateDB()
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err, "db should be created successfully.")

	// defer tearDown(t, dbMigration)

	assert.Nil(t, dbMigration.MigrateDB(), "db should be migrate successfully")
}
