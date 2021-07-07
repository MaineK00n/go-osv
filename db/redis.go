package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MaineK00n/go-osv/config"
	"github.com/MaineK00n/go-osv/models"
	"github.com/go-redis/redis/v8"
	"github.com/inconshreveable/log15"
	"golang.org/x/xerrors"
	pb "gopkg.in/cheggaaa/pb.v1"
)

/**
# Redis Data Structure
- HASH
  ┌───┬──────────────┬────────────────────────────────────────┬──────────┬───────────────────────────┐
  │NO │     HASH     │               FIELD                    │  VALUE   │          PURPOSE          │
  └───┴──────────────┴────────────────────────────────────────┴──────────┴───────────────────────────┘
  ┌───┬──────────────┬────────────────────────────────────────┬──────────┬───────────────────────────┐
  │ 1 │OSV#$PKGNAME  │  crates.io/DWF/Go/Linux/OSS-Fuzz/PyPI  │ $CVEJSON │  TO GET JSON BY PKG Name  │
  └───┴──────────────┴────────────────────────────────────────┴──────────┴───────────────────────────┘
  ┌───┬──────────────┬────────────────────────────────────────┬──────────┬───────────────────────────┐
  │ 2 │   OSV#$ID    │  crates.io/DWF/Go/Linux/OSS-Fuzz/PyPI  │ $CVEJSON │     TO GET JSON BY ID     │
  └───┴──────────────┴────────────────────────────────────────┴──────────┴───────────────────────────┘
**/

const (
	dialectRedis  = "redis"
	hashKeyPrefix = "OSV#"
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

// IsGoOSVModelV1 determines if the DB was created at the time of go-osv Model v1
func (r *RedisDriver) IsGoOSVModelV1() (bool, error) {
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
func (r *RedisDriver) InsertOSVs(osvType models.OSVType, osvJSONs []models.OSVJSON) error {
	ctx := context.Background()
	osvs, err := models.ConvertOSV(osvJSONs)
	if err != nil {
		return err
	}
	bar := pb.StartNew(len(osvs))

	for _, osv := range osvs {
		pipe := r.conn.Pipeline()
		bar.Increment()

		j, err := json.Marshal(osv)
		if err != nil {
			return fmt.Errorf("Failed to marshal json. err: %s", err)
		}

		if result := pipe.HSet(ctx, hashKeyPrefix+osv.Package.Name, string(osvType), string(j)); result.Err() != nil {
			return fmt.Errorf("Failed to HSet CVE. err: %s", result.Err())
		}

		for _, aliase := range osv.Aliases {
			if result := pipe.HSet(ctx, hashKeyPrefix+aliase.Alias, string(osvType), string(j)); result.Err() != nil {
				return fmt.Errorf("Failed to HSet CVE. err: %s", result.Err())
			}
		}

		if _, err = pipe.Exec(ctx); err != nil {
			return fmt.Errorf("Failed to exec pipeline. err: %s", err)
		}
	}
	bar.Finish()

	return nil
}

// GetOSVbyID :
func (r *RedisDriver) GetOSVbyID(ID string, osvType string) ([]models.OSV, error) {
	ctx := context.Background()
	result := r.conn.HGetAll(ctx, hashKeyPrefix+ID)
	if result.Err() != nil {
		return nil, result.Err()
	}

	osvs := []models.OSV{}
	if osvType != "" {
		if j, ok := result.Val()[osvType]; ok {
			osv := models.OSV{}
			if err := json.Unmarshal([]byte(j), &osv); err != nil {
				return nil, err
			}
			osvs = append(osvs, osv)
		}
	} else {
		for _, j := range result.Val() {
			osv := models.OSV{}
			if err := json.Unmarshal([]byte(j), &osv); err != nil {
				return nil, err
			}
			osvs = append(osvs, osv)
		}
	}

	return osvs, nil
}

// GetOSVbyPackageName :
func (r *RedisDriver) GetOSVbyPackageName(name string, osvType string) ([]models.OSV, error) {
	ctx := context.Background()

	result := r.conn.HGetAll(ctx, hashKeyPrefix+name)
	if result.Err() != nil {
		return nil, result.Err()
	}

	osvs := []models.OSV{}
	if osvType != "" {
		if j, ok := result.Val()[osvType]; ok {
			osv := models.OSV{}
			if err := json.Unmarshal([]byte(j), &osv); err != nil {
				return nil, err
			}
			osvs = append(osvs, osv)
		}
	} else {
		for _, j := range result.Val() {
			osv := models.OSV{}
			if err := json.Unmarshal([]byte(j), &osv); err != nil {
				return nil, err
			}
			osvs = append(osvs, osv)
		}
	}

	return osvs, nil
}
