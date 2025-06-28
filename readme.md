# TitanLift Back End

Golang backend server designed for an application that logs gym progress.


## Requirements
* Go 1.24+
* PostgreSQL Database

## Configuration

The server is configured using a TOML configuration file. Template is available in the `/configs` folder.

The database connection is configured in the following environment variable:

- `TITAN_DB_URL` - URL of PostgresDB connection (ex: postgres://username:password@127.0.0.1:5432/database_name?sslmode=disable)