package migration

import (
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

const (
	dbPath                 = "testdata/test_data.db"
	testUsersTableName     = "migrate_users_test"
	testAddressesTableName = "migrate_addresses_test"
)

func newTestMigrateDB(db *sqlx.DB) (*DBMigration, error) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	// open bolt db for migration
	boltDB, err := bolt.Open(dbPath, 0600, nil)
	if boltDB == nil {
		return nil, err
	}

	// initiate postgres for migration
	var schema = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS "%[1]s" (
  id    SERIAL PRIMARY KEY,
  email text NOT NULL UNIQUE,
  last_updated TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS "%[2]s" (
  id        SERIAL PRIMARY KEY,
  address   text      NOT NULL UNIQUE,
  timestamp TIMESTAMP NOT NULL DEFAULT now(),
  user_id   SERIAL    NOT NULL REFERENCES "%[1]s" (id)
);
`, testUsersTableName, testAddressesTableName)

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	defer pgsql.CommitOrRollback(tx, sugar, &err)

	if _, err = tx.Exec(schema); err != nil {
		return nil, err
	}

	return &DBMigration{sugar, boltDB, db}, nil
}

func TestMigrateDB(t *testing.T) {
	db, teardown := testutil.MustNewDevelopmentDB()
	dbMigration, err := newTestMigrateDB(db)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err, "db should be created successfully.")

	defer func() {
		assert.NoError(t, teardown())
	}()

	assert.Nil(t, dbMigration.Migrate(testUsersTableName, testAddressesTableName), "db should be migrate successfully")
}
