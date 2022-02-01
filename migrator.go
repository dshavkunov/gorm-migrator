package gorm_migrator

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sulimak-co/gorm-migrator/internal/helpers"
	"github.com/sulimak-co/gorm-migrator/internal/migration"
	"github.com/sulimak-co/gorm-migrator/internal/migration/repository"
	"github.com/sulimak-co/gorm-migrator/internal/models"
	"gorm.io/gorm"
)

type Migration struct {
	Up   func(migrator gorm.Migrator) error
	Down func(migrator gorm.Migrator) error
}

// migrator represents migration manager.
type migrator struct {
	logger        logrus.FieldLogger
	migrator      gorm.Migrator
	migrationRepo migration.Repository
	allMigrations []Migration
}

type Migrator interface {
	Execute() error
}

// NewMigrator creates new migration manager.
func NewMigrator(conn *gorm.DB, allMigrations []Migration) Migrator {
	return &migrator{
		logger:        logrus.New(),
		migrator:      conn.Migrator(),
		migrationRepo: repository.NewMigrationRepository(conn),
		allMigrations: allMigrations,
	}
}

// Execute applies new database migrations.
func (m *migrator) Execute() error {
	if err := createMigrationTable(m.migrator); err != nil {
		return helpers.WrapError(err, "migration init error")
	}

	if err := m.validateNewMigrations(); err != nil {
		return helpers.WrapError(err, "validation error")
	}

	_, err := m.up()
	if err != nil {
		return helpers.WrapError(err, "migrate up error")
	}

	return nil
}

func (m *migrator) validateNewMigrations() error {
	newIndexes, err := m.up()
	if err != nil {
		return helpers.WrapError(err, "validation error")
	}
	if err := m.down(newIndexes); err != nil {
		return helpers.WrapError(err, "validation error")
	}
	return nil
}

func (m *migrator) up() ([]int, error) {
	var newMigrationIndexes []int
	for i, mgrt := range m.allMigrations {
		migrationName, err := helpers.GetFunctionName(mgrt.Up)
		if err != nil {
			return nil, err
		}
		mRecord, err := m.migrationRepo.GetByName(migrationName)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
		}
		if mRecord != nil {
			continue
		}
		mRecord = &models.Migration{}
		newMigrationIndexes = append(newMigrationIndexes, i)
		if err := mgrt.Up(m.migrator); err != nil {
			if downErr := m.down(newMigrationIndexes); downErr != nil {
				m.logger.Errorf("migration down error: %v", err)
			}
			return nil, err
		}
		mRecord.Name = migrationName
		if err := m.migrationRepo.Store(mRecord); err != nil {
			return nil, err
		}
	}
	return newMigrationIndexes, nil
}

// down rollbacks migrations.
func (m *migrator) down(newMigrationIndexes []int) error {
	count := len(newMigrationIndexes)
	for i := count - 1; i >= 0; i-- {
		mgrt := m.allMigrations[newMigrationIndexes[i]]
		if err := mgrt.Down(m.migrator); err != nil {
			return err
		}
		migrationName, err := helpers.GetFunctionName(mgrt.Up)
		if err != nil {
			return err
		}
		if err := m.migrationRepo.DeleteByName(migrationName); err != nil {
			return err
		}
	}
	return nil
}

// createMigrationTable creates migration table.
func createMigrationTable(migrator gorm.Migrator) error {
	if migrator.HasTable(&models.Migration{}) {
		return nil
	}

	return migrator.CreateTable(&models.Migration{})
}
