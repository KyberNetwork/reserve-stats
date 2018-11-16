package migration

import (
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
	_, err := dbm.postgresDB.Exec(fmt.Sprintf(`DROP TABLE "%s"`, tableName))
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
	CREATE TABLE IF NOT EXISTS "migrate_price_analytics" (
	  id SERIAL PRIMARY KEY,
	  timestamp TIMESTAMP NOT NULL,
	  block_expiration BOOL NOT NULL
	);
	CREATE TABLE IF NOT EXISTS "migrate_price_analytic_data" (
		id SERIAL PRIMARY KEY,
		token text NOT NULL,
		ask_price float NOT NULL,
		bid_price float NOT NULL,
		mid_afp_price float NOT NULL,
		mid_afp_old_price float NOT NULL,
		min_spread float NOT NULL,
		trigger_update bool NOT NULL,
		price_analytic_id SERIAL NOT NULL REFERENCES migrate_price_analytics (id)
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
	assert.Nil(t, dbMigration.deleteTable("migrate_price_analytic_data"), "test table should be tear down succesfully.")
	assert.Nil(t, dbMigration.deleteTable("migrate_price_analytics"), "test table should be tear down succesfully.")
}

func TestMigrateDB(t *testing.T) {
	t.Skip()
	dbMigration, err := newTestMigrateDB()
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err, "db should be created successfully.")

	defer tearDown(t, dbMigration)

	assert.Nil(t, dbMigration.Migrate(), "db should be migrate successfully")
}
