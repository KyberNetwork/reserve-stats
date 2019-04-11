package testutil

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // sql driver name: "postgres"
)

const (
	postgresHost     = "127.0.0.1"
	postgresPort     = 5432
	postgresUser     = "reserve_stats"
	postgresPassword = "reserve_stats"
	postgresDatabase = "reserve_stats"
)

// MustNewDevelopmentDB returns a new development db instance.
func MustNewDevelopmentDB() (dbName string, db *sqlx.DB) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
	)
	db = sqlx.MustConnect("postgres", connStr)
	return postgresDatabase, db
}

// MustNewRandomDevelopmentDB creates a new development DB.
// It also returns a function to teardown it after the test.
func MustNewRandomDevelopmentDB() (ddlDB *sqlx.DB, teardown func() error) {
	dbName := RandomString(8)

	ddlDBConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresDatabase,
	)
	ddlDB = sqlx.MustConnect("postgres", ddlDBConnStr)
	ddlDB.MustExec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	if err := ddlDB.Close(); err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		dbName,
	)
	db := sqlx.MustConnect("postgres", connStr)
	return db, func() error {
		if err := db.Close(); err != nil {
			return err
		}
		ddlDB, err := sqlx.Connect("postgres", ddlDBConnStr)
		if err != nil {
			return err
		}

		if _, err = ddlDB.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbName)); err != nil {
			return err
		}

		return ddlDB.Close()
	}
}
