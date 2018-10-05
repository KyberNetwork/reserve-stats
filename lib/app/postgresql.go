package app

import (
	"github.com/go-pg/pg"
	"github.com/urfave/cli"
)

const (
	postgresHostFlag     = "postgres_host"
	postgresUserFlag     = "postgres_user"
	postgresPasswordFlag = "postgres_password"
	postgresDatabaseFlag = "postgres_database"
)

// NewPostgreSQLFlags creates new cli flags for PostgreSQL client.
func NewPostgreSQLFlags(defaultDB string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   postgresHostFlag,
			Usage:  "PostgreSQL host to connect",
			EnvVar: joinEnvVar(commonEnvPrefix, "POSTGRES_HOST"),
			Value:  "127.0.0.1:5432",
		},
		cli.StringFlag{
			Name:   postgresUserFlag,
			Usage:  "PostgreSQL user to connect",
			EnvVar: joinEnvVar(commonEnvPrefix, "POSTGRES_USER"),
			Value:  "",
		},
		cli.StringFlag{
			Name:   postgresPasswordFlag,
			Usage:  "PostgreSQL password to connect",
			EnvVar: joinEnvVar(commonEnvPrefix, "POSTGRES_PASSWORD"),
			Value:  "",
		},
		cli.StringFlag{
			Name:   postgresDatabaseFlag,
			Usage:  "Postgres database to connect",
			EnvVar: joinEnvVar(commonEnvPrefix, "POSTGRES_DATABASE"),
			Value:  defaultDB,
		},
	}
}

// NewDBFromContext creates a DB instance from cli flags configuration.
func NewDBFromContext(c *cli.Context) *pg.DB {
	return pg.Connect(&pg.Options{
		Addr:     c.String(postgresHostFlag),
		User:     c.String(postgresUserFlag),
		Password: c.String(postgresPasswordFlag),
		Database: c.String(postgresDatabaseFlag),
	})
}
