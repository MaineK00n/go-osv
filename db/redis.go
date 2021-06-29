package db

import (
	"context"
	"fmt"

	"github.com/MaineK00n/go-osv/config"
	"github.com/MaineK00n/go-osv/models"
	"github.com/go-redis/redis/v8"
	"github.com/inconshreveable/log15"
	"golang.org/x/xerrors"
)

/**
# Redis Data Structure
**/

const (
	dialectRedis = "redis"
)

// RedisDriver is Driver for Redis
type RedisDriver struct {
	name string
	conn *redis.Client
}

// Name return db name
func (r *RedisDriver) Name() string {
	return r.name
}

// OpenDB opens Database
func (r *RedisDriver) OpenDB(dbType, dbPath string, debugSQL bool) (locked bool, err error) {
	if err = r.connectRedis(dbPath); err != nil {
		err = fmt.Errorf("Failed to open DB. dbtype: %s, dbpath: %s, err: %s", dbType, dbPath, err)
	}
	return
}

func (r *RedisDriver) connectRedis(dbPath string) error {
	ctx := context.Background()
	var err error
	var option *redis.Options
	if option, err = redis.ParseURL(dbPath); err != nil {
		log15.Error("Failed to parse url.", "err", err)
		return err
	}
	r.conn = redis.NewClient(option)
	err = r.conn.Ping(ctx).Err()
	return err
}

// CloseDB close Database
func (r *RedisDriver) CloseDB() (err error) {
	if r.conn == nil {
		return
	}
	if err = r.conn.Close(); err != nil {
		return xerrors.Errorf("Failed to close DB. Type: %s. err: %w", r.name, err)
	}
	return
}

// MigrateDB migrates Database
func (r *RedisDriver) MigrateDB() error {
	return nil
}

// IsGostModelV1 determines if the DB was created at the time of Gost Model v1
func (r *RedisDriver) IsGostModelV1() (bool, error) {
	return false, nil
}

// GetFetchMeta get FetchMeta from Database
func (r *RedisDriver) GetFetchMeta() (*models.FetchMeta, error) {
	return &models.FetchMeta{GoOSVRevision: config.Revision, SchemaVersion: models.LatestSchemaVersion}, nil
}

// UpsertFetchMeta upsert FetchMeta to Database
func (r *RedisDriver) UpsertFetchMeta(*models.FetchMeta) error {
	return nil
}

// InsertOSVs :
func (r *RedisDriver) InsertOSVs(_ models.OSVType, osvs []models.OSVJSON) error {
	return nil
}
