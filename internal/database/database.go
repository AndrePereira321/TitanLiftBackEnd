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

func New(serverConfig *config.ServerConfig) (*Database, error) {
	dbLogger, err := logger.New("DATABASE", serverConfig.Logging().DatabaseLogLevel(), serverConfig.Logging().LogDir())
	if err != nil {
		return nil, err
	}

	dbURL := os.Getenv(config.EnvTitanDbUrl)
	if dbURL == "" {
		return nil, server_error.New("DB_UPGRADE", "database URL environment variable is not set")
	}

	dbLogger.Info("Initializing database.")
	err = upgradeStructure(serverConfig, dbLogger, dbURL)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, server_error.Wrap("DB_CONNECT", "failed to open database", err)
	}

	db.SetMaxIdleConns(serverConfig.Database().MaxIdleCons())
	db.SetMaxOpenConns(serverConfig.Database().MaxOpenCons())

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, server_error.Wrap("DB_CONNECT", "failed to ping database", err)
	}

	dbLogger.Info("Database successfully initialized.")

	return &Database{
		db:     db,
		logger: dbLogger,
	}, nil
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

func upgradeStructure(serverConfig *config.ServerConfig, dbLogger *logger.Logger, dbURL string) error {
	dbLogger.Debug("Upgrading database structure")

	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "failed to parse database URL", err)
	}

	dbMateLogFile, err := url.JoinPath(serverConfig.Logging().LogDir(), "DB_UPGRADE.log")
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "failed to create DBMate log file path\n", err)
	}

	dbMateFile, err := os.OpenFile(dbMateLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "failed to open DBMate log file", err)
	}

	defer func(dbMateFile *os.File) {
		_ = dbMateFile.Close()
	}(dbMateFile)

	dbMate := dbmate.New(parsedURL)
	dbMate.FS = fs
	dbMate.Verbose = true
	dbMate.Strict = true
	dbMate.Log = dbMateFile
	dbMate.MigrationsDir = []string{"migrations"}

	migrations, err := dbMate.FindMigrations()
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "failed to find migrations", err)
	}

	for _, migration := range migrations {
		dbLogger.Trace(fmt.Sprintf("Found migration %s with version %s. Applied: %t", migration.FileName, migration.Version, migration.Applied))
	}

	dbLogger.Trace("Applying migrations.")
	err = dbMate.Migrate()
	if err != nil {
		return server_error.Wrap("DB_UPGRADE", "failed to apply migrations", err)
	}

	dbLogger.Debug("Database structure upgrade completed successfully")
	return nil
}
