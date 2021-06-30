package db

import (
	"fmt"

	"github.com/MaineK00n/go-osv/models"
	"github.com/inconshreveable/log15"
	"golang.org/x/xerrors"
)

// DB is interface for a database driver
type DB interface {
	Name() string
	OpenDB(string, string, bool) (bool, error)
	CloseDB() error
	MigrateDB() error

	IsGostModelV1() (bool, error)
	GetFetchMeta() (*models.FetchMeta, error)
	UpsertFetchMeta(*models.FetchMeta) error

	InsertOSVs(models.OSVType, []models.OSVJSON) error
	GetOSVbyID(ID string, osvType string) ([]models.OSV, error)
	GetOSVbyPackageName(name string, osvType string) ([]models.OSV, error)
}

// NewDB returns db driver
func NewDB(dbType, dbPath string, debugSQL bool) (driver DB, locked bool, err error) {
	if driver, err = newDB(dbType); err != nil {
		log15.Error("Failed to new db.", "err", err)
		return driver, false, err
	}

	if locked, err := driver.OpenDB(dbType, dbPath, debugSQL); err != nil {
		if locked {
			return nil, true, err
		}
		return nil, false, err
	}

	isV1, err := driver.IsGostModelV1()
	if err != nil {
		log15.Error("Failed to IsGostModelV1.", "err", err)
		return nil, false, err
	}
	if isV1 {
		log15.Error("Failed to NewDB. Since SchemaVersion is incompatible, delete Database and fetch again")
		return nil, false, xerrors.New("Failed to NewDB. Since SchemaVersion is incompatible, delete Database and fetch again.")
	}

	if err := driver.MigrateDB(); err != nil {
		log15.Error("Failed to migrate db.", "err", err)
		return driver, false, err
	}
	return driver, false, nil
}

func newDB(dbType string) (DB, error) {
	switch dbType {
	case dialectSqlite3, dialectMysql, dialectPostgreSQL:
		return &RDBDriver{name: dbType}, nil
	case dialectRedis:
		return &RedisDriver{name: dbType}, nil
	}
	return nil, fmt.Errorf("Invalid database dialect. err: %s", dbType)
}
