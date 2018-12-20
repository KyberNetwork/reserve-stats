package migration

import (
	"fmt"
	"testing"

	"go.uber.org/zap"

	"github.com/boltdb/bolt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
)

const (
	postgresHost           = "127.0.0.1"
	postgresPort           = 5432
	postgresUser           = "reserve_stats"
	postgresPassword       = "reserve_stats"
	postgresDatabase       = "reserve_stats"
	dbPath                 = "testdata/test_data.db"
	testUsersTableName     = "migrate_users_test"
	testAddressesTableName = "migrate_addresses_test"
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
	var schema = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS "%[1]s" (
  id    SERIAL PRIMARY KEY,
  email text NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS "%[2]s" (
  id        SERIAL PRIMARY KEY,
  address   text      NOT NULL UNIQUE,
  timestamp TIMESTAMP NOT NULL DEFAULT now(),
  user_id   SERIAL    NOT NULL REFERENCES "%[1]s" (id)
);
`, testUsersTableName, testAddressesTableName)

	tx, err := postgres.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback()
			if rollBackErr != err {
				logger.Debug(fmt.Sprintf("rollback error: %s", rollBackErr))
			}
			return
		}
		err = tx.Commit()
	}()

	if _, err = tx.Exec(schema); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &DBMigration{sugar, boltDB, postgres}, nil
}

func tearDown(t *testing.T, dbMigration *DBMigration) {
	assert.Nil(t, dbMigration.deleteTable(testAddressesTableName), "test table should be tear down succesfully.")
	assert.Nil(t, dbMigration.deleteTable(testUsersTableName), "test table should be tear down succesfully.")
}

func TestMigrateDB(t *testing.T) {
	dbMigration, err := newTestMigrateDB()
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err, "db should be created successfully.")

	defer tearDown(t, dbMigration)

	assert.Nil(t, dbMigration.Migrate(testUsersTableName, testAddressesTableName), "db should be migrate successfully")
}
