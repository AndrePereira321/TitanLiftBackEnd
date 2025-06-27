# TitanLift Back End

Golang backend server designed for an application that logs gym progress.


## Requirements
* Go 1.24+
* PostgreSQL Database

## Configuration

The server is configured using a TOML configuration file. Template is available in the `/configs` folder.

The database connection is configured in the following environment variables:

- `TITAN_DB_HOST` - IP address of the database server
- `TITAN_DB_USER` - Username for the database
- `TITAN_DB_PASSWORD` - Password for the database user
- `TITAN_DB_NAME` - Name of the database