package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MaineK00n/go-osv/config"
	"github.com/MaineK00n/go-osv/models"
	"github.com/inconshreveable/log15"
	sqlite3 "github.com/mattn/go-sqlite3"
	"golang.org/x/xerrors"
	pb "gopkg.in/cheggaaa/pb.v1"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// Supported DB dialects.
const (
	dialectSqlite3    = "sqlite3"
	dialectMysql      = "mysql"
	dialectPostgreSQL = "postgres"
)

// RDBDriver is Driver for RDB
type RDBDriver struct {
	name string
	conn *gorm.DB
}

// Name return db name
func (r *RDBDriver) Name() string {
	return r.name
}

// OpenDB opens Database
func (r *RDBDriver) OpenDB(dbType, dbPath string, debugSQL bool) (locked bool, err error) {
	gormConfig := gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
	}

	if debugSQL {
		gormConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	}

	switch r.name {
	case dialectSqlite3:
		r.conn, err = gorm.Open(sqlite.Open(dbPath), &gormConfig)
	case dialectMysql:
		r.conn, err = gorm.Open(mysql.Open(dbPath), &gormConfig)
	case dialectPostgreSQL:
		r.conn, err = gorm.Open(postgres.Open(dbPath), &gormConfig)
	default:
		err = xerrors.Errorf("Not Supported DB dialects. r.name: %s", r.name)
	}

	if err != nil {
		msg := fmt.Sprintf("Failed to open DB. dbtype: %s, dbpath: %s, err: %s", dbType, dbPath, err)
		if r.name == dialectSqlite3 {
			switch err.(sqlite3.Error).Code {
			case sqlite3.ErrLocked, sqlite3.ErrBusy:
				return true, fmt.Errorf(msg)
			}
		}
		return false, fmt.Errorf(msg)
	}

	if r.name == dialectSqlite3 {
		r.conn.Exec("PRAGMA foreign_keys = ON")
	}
	return false, nil
}

// CloseDB close Database
func (r *RDBDriver) CloseDB() (err error) {
	if r.conn == nil {
		return
	}

	var sqlDB *sql.DB
	if sqlDB, err = r.conn.DB(); err != nil {
		return xerrors.Errorf("Failed to get DB Object. err : %w", err)
	}
	if err = sqlDB.Close(); err != nil {
		return xerrors.Errorf("Failed to close DB. Type: %s. err: %w", r.name, err)
	}
	return
}

// MigrateDB migrates Database
func (r *RDBDriver) MigrateDB() error {
	if err := r.conn.AutoMigrate(
		&models.FetchMeta{},

		&models.OSV{},
		&models.OSVAliases{},
		&models.OSVRelated{},
		&models.OSVPackage{},
		&models.OSVAffects{},
		&models.OSVAffectsRanges{},
		&models.OSVAffectsVersions{},
		&models.OSVReferences{},
		&models.OSVEcosystemSpecific{},
		&models.OSVDatabaseSpecific{},
	); err != nil {
		return xerrors.Errorf("Failed to migrate. err: %w", err)
	}

	return nil
}

// IsGostModelV1 determines if the DB was created at the time of Gost Model v1
func (r *RDBDriver) IsGostModelV1() (bool, error) {
	if r.conn.Migrator().HasTable(&models.FetchMeta{}) {
		return false, nil
	}

	var (
		count int64
		err   error
	)
	switch r.name {
	case dialectSqlite3:
		err = r.conn.Table("sqlite_master").Where("type = ?", "table").Count(&count).Error
	case dialectMysql:
		err = r.conn.Table("information_schema.tables").Where("table_schema = ?", r.conn.Migrator().CurrentDatabase()).Count(&count).Error
	case dialectPostgreSQL:
		err = r.conn.Table("pg_tables").Where("schemaname = ?", "public").Count(&count).Error
	}

	if count > 0 {
		return true, nil
	}
	return false, err
}

// GetFetchMeta get FetchMeta from Database
func (r *RDBDriver) GetFetchMeta() (fetchMeta *models.FetchMeta, err error) {
	if err = r.conn.Take(&fetchMeta).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return &models.FetchMeta{GoOSVRevision: config.Revision, SchemaVersion: models.LatestSchemaVersion}, nil
	}

	return fetchMeta, nil
}

// UpsertFetchMeta upsert FetchMeta to Database
func (r *RDBDriver) UpsertFetchMeta(fetchMeta *models.FetchMeta) error {
	fetchMeta.GoOSVRevision = config.Revision
	fetchMeta.SchemaVersion = models.LatestSchemaVersion
	return r.conn.Save(fetchMeta).Error
}

// IndexChunk has a starting point and an ending point for Chunk
type IndexChunk struct {
	From, To int
}

func chunkSlice(length int, chunkSize int) <-chan IndexChunk {
	ch := make(chan IndexChunk)

	go func() {
		defer close(ch)

		for i := 0; i < length; i += chunkSize {
			idx := IndexChunk{i, i + chunkSize}
			if length < idx.To {
				idx.To = length
			}
			ch <- idx
		}
	}()

	return ch
}

// InsertOSVs :
func (r *RDBDriver) InsertOSVs(insertOSVType models.OSVType, osvJSONs []models.OSVJSON) error {
	osvs, err := models.ConvertOSV(osvJSONs)
	if err != nil {
		return xerrors.Errorf("Failed to Convert OSV JSON. er: %w", err)
	}

	if err := r.deleteAndInsertOSVs(r.conn, insertOSVType, osvs); err != nil {
		return xerrors.Errorf("Failed to insert OSV data. err: %w", err)
	}
	return nil
}

func (r *RDBDriver) deleteAndInsertOSVs(conn *gorm.DB, insertOSVType models.OSVType, osvs []models.OSV) (err error) {
	bar := pb.StartNew(len(osvs))
	tx := conn.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	// Delete all old records
	oldIDs := []int64{}
	result := tx.Model(models.OSVPackage{}).Select("osv_id").Where("ecosystem = ?", insertOSVType).Find(&oldIDs)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return xerrors.Errorf("Failed to select old defs: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		affectsIDs := []int64{}
		result := tx.Model(models.OSVAffects{}).Select("id").Where("osv_id IN ?", oldIDs).Find(&affectsIDs)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return xerrors.Errorf("Failed to select old affects: %w", result.Error)
		}

		if result.RowsAffected > 0 {
			if err := tx.Unscoped().Where("osv_affects_id IN ?", affectsIDs).Delete(&models.OSVAffectsRanges{}).Error; err != nil {
				tx.Rollback()
				return xerrors.Errorf("Failed to delete: %w", err)
			}

			if err := tx.Unscoped().Where("osv_affects_id IN ?", affectsIDs).Delete(&models.OSVAffectsVersions{}).Error; err != nil {
				tx.Rollback()
				return xerrors.Errorf("Failed to delete: %w", err)
			}
		}

		if err := tx.Unscoped().Where("osv_id IN ?", oldIDs).Delete(&models.OSVAffects{}).Error; err != nil {
			tx.Rollback()
			return xerrors.Errorf("Failed to delete: %w", err)
		}

		if err := tx.Unscoped().Where("osv_id IN ?", oldIDs).Delete(&models.OSVAliases{}).Error; err != nil {
			tx.Rollback()
			return xerrors.Errorf("Failed to delete: %w", err)
		}

		if err := tx.Unscoped().Where("osv_id IN ?", oldIDs).Delete(&models.OSVRelated{}).Error; err != nil {
			tx.Rollback()
			return xerrors.Errorf("Failed to delete: %w", err)
		}

		if err := tx.Unscoped().Where("osv_id IN ?", oldIDs).Delete(&models.OSVPackage{}).Error; err != nil {
			tx.Rollback()
			return xerrors.Errorf("Failed to delete: %w", err)
		}

		if err := tx.Unscoped().Where("osv_id IN ?", oldIDs).Delete(&models.OSVReferences{}).Error; err != nil {
			tx.Rollback()
			return xerrors.Errorf("Failed to delete: %w", err)
		}

		if err := tx.Unscoped().Where("osv_id IN ?", oldIDs).Delete(&models.OSVEcosystemSpecific{}).Error; err != nil {
			tx.Rollback()
			return xerrors.Errorf("Failed to delete: %w", err)
		}

		if err := tx.Unscoped().Where("osv_id IN ?", oldIDs).Delete(&models.OSVDatabaseSpecific{}).Error; err != nil {
			tx.Rollback()
			return xerrors.Errorf("Failed to delete: %w", err)
		}

		if err := tx.Unscoped().Where("id IN ?", oldIDs).Delete(&models.OSV{}).Error; err != nil {
			tx.Rollback()
			return xerrors.Errorf("Failed to delete: %w", err)
		}

	}

	for idx := range chunkSlice(len(osvs), 10) {
		if err = tx.Create(osvs[idx.From:idx.To]).Error; err != nil {
			return fmt.Errorf("Failed to insert. err: %s", err)
		}
		bar.Add(idx.To - idx.From)
	}
	bar.Finish()

	return nil
}

// GetOSVbyID :
func (r *RDBDriver) GetOSVbyID(ID string, osvType string) ([]models.OSV, error) {
	osvIDs := []int64{}

	q := r.conn.Model(models.OSVAliases{}).Select("osv_aliases.osv_id")
	if osvType != "" {
		q = q.Joins("JOIN osv_packages ON osv_aliases.osv_id = osv_packages.osv_id").Where("ecosystem = ? AND alias = ?", osvType, ID)
	} else {
		q = q.Where("alias = ?", ID)
	}

	if err := q.Find(&osvIDs).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log15.Error("Failed to GetOSVbyID", "err", err)
			return []models.OSV{}, err
		}
		return []models.OSV{}, nil
	}

	osv := []models.OSV{}
	if err := r.conn.Preload("Affects.Ranges").Preload("Affects.Versions").Preload(clause.Associations).Where("id IN ?", osvIDs).Find(&osv).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log15.Error("Failed to GetOSVbyID", "err", err)
			return []models.OSV{}, err
		}
		return []models.OSV{}, nil
	}

	return osv, nil
}

// GetOSVbyPackageName :
func (r *RDBDriver) GetOSVbyPackageName(name string, osvType string) ([]models.OSV, error) {
	osvIDs := []int64{}

	q := r.conn.Model(models.OSVPackage{}).Select("osv_packages.osv_id")
	if osvType != "" {
		q = q.Where("ecosystem = ? AND name = ?", osvType, name)
	} else {
		q = q.Where("name = ?", name)
	}

	if err := q.Find(&osvIDs).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.OSV{}, err
		}
		return []models.OSV{}, nil
	}

	osv := []models.OSV{}
	if err := r.conn.Preload("Affects.Ranges").Preload("Affects.Versions").Preload(clause.Associations).Where("id IN ?", osvIDs).Find(&osv).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.OSV{}, err
		}
		return []models.OSV{}, nil
	}

	return osv, nil
}
