package database

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	_ "github.com/lib/pq"
	"net/url"
	"os"
	"titan-lift/internal/config"
	"titan-lift/internal/logger"
	"titan-lift/internal/server_error"
)

type Database struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewDatabase(config *config.ServerConfig) (*sql.DB, error) {
	dbLogger, err := logger.NewLogger("DATABASE", config.Logging().DatabaseLogLevel(), config.Logging().LogDir())
	if err != nil {
		return nil, err
	}

	dbLogger.Info("Initializing database.")
	err = upgradeStructure(config, dbLogger)
	if err != nil {
		return nil, err
	}

	dbLogger.Info("Database successfully initialized.")
	return nil, nil
}

func (d *Database) Close() error {
	err := d.db.Close()
	if err != nil {
		return server_error.Wrap("DB_CLOSE", "error when closing database", err)
	}
	return d.logger.Close()
}

//go:embed migrations/*.sql
var fs embed.FS

func upgradeStructure(serverConfig *config.ServerConfig, dbLogger *logger.Logger) error {
	dbLogger.Debug("Upgrading database structure.")

	dbUrl := os.Getenv(config.EnvTitanDbUrl)
	if len(dbUrl) == 0 {
		return server_error.New("DB_UPGRADE", "unable to find database url")
	}

	u, err := url.Parse(dbUrl)
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "error when parsing database url", err)
	}

	dbMateLogFile, err := url.JoinPath(serverConfig.Logging().LogDir(), "DB_UPGRADE.log")
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "error when joining dbmate log file path", err)
	}

	dbMateFile, err := os.OpenFile(dbMateLogFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "error when opening dbmate log file", err)
	}

	defer func(dbMateFile *os.File) {
		_ = dbMateFile.Close()
	}(dbMateFile)

	dbMate := dbmate.New(u)
	dbMate.FS = fs
	dbMate.Verbose = true
	dbMate.Strict = true
	dbMate.Log = dbMateFile
	dbMate.MigrationsDir = []string{"migrations"}

	migrations, err := dbMate.FindMigrations()
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "error when finding migrations", err)
	}

	for _, migration := range migrations {
		dbLogger.Trace(fmt.Sprintf("Found migration %s with version %s. Applied: %t", migration.FileName, migration.Version, migration.Applied))
	}

	dbLogger.Trace("Applying migrations.")
	err = dbMate.Migrate()
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "error when upgrading database", err)
	}

	dbLogger.Debug("Database structure successfully upgraded.")
	return nil
}
