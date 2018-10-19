package migration

import (
	"fmt"
	"go.uber.org/zap"
	"testing"

	"github.com/boltdb/bolt"
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
	dbPath           = "testdata/test_data.db"
)

//DeleteTable delete a table from  schema
func (dbm *DBMigration) deleteTable(tableName string) error {
	_, err := dbm.postgresdb.Exec(fmt.Sprintf(`DROP TABLE "%s"`, tableName))
	return err
}

func newTestMigrateDB() (*DBMigration, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
	)
	postgres, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// open bolt db for migration
	boltDB, err := bolt.Open(dbPath, 0600, nil)
	if boltDB == nil {
		return nil, err
	}

	// initiate postgres for migration
	const schema = `
	CREATE TABLE IF NOT EXISTS "migrate_users_test" (
	  id    SERIAL PRIMARY KEY,
	  email text NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS "migrate_addresses_test" (
	  id        SERIAL PRIMARY KEY,
	  address   text      NOT NULL UNIQUE,
	  timestamp TIMESTAMP NOT NULL DEFAULT now(),
	  user_id   SERIAL    NOT NULL REFERENCES migrate_users_test (id)
	);
	`

	tx, err := postgres.Beginx()
	if err != nil {
		return nil, err
	}

	if _, err = tx.Exec(schema); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &DBMigration{sugar, boltDB, postgres}, nil
}

func tearDown(t *testing.T, dbMigration *DBMigration) {
	assert.Nil(t, dbMigration.deleteTable("migrate_addresses_test"), "test table should be tear down succesfully.")
	assert.Nil(t, dbMigration.deleteTable("migrate_users_test"), "test table should be tear down succesfully.")
}

func TestMigrateDB(t *testing.T) {
	dbMigration, err := newTestMigrateDB()
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err, "db should be created successfully.")

	defer tearDown(t, dbMigration)

	assert.Nil(t, dbMigration.Migrate(), "db should be migrate successfully")
}
