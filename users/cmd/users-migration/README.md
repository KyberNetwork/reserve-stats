# Migrate
Migrate database user from bolt to postgresql

## Build

```shell
cd cmd
go build -v .
```

## Run

```
./cmd --postgres-user=reserve_stats --postgres-password=reserve_stats --postgres-database=reserve_stats --database=dev_users.db
```

Available flags:

**postgres-user**: user to connect to postgres
**postgres-password**: password for postgres database
**postgres-database**: database to migrate db to 
**database**: path to boltdb data file
