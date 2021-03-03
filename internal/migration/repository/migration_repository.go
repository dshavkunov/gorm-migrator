// Package repository describes migration database methods.
package repository

import (
	"github.com/sulimak-co/gorm-migrator/internal/helpers"
	"github.com/sulimak-co/gorm-migrator/internal/migration"
	"github.com/sulimak-co/gorm-migrator/internal/models"
	"gorm.io/gorm"
)

const tableName = "migrations"

// migrationRepository defines a migration repository.
type migrationRepository struct {
	Conn *gorm.DB
}

// NewMigrationRepository creates new interface of migration repository.
func NewMigrationRepository(conn *gorm.DB) migration.Repository {
	return &migrationRepository{conn}
}

// GetByName returns the migration that matches the given name.
func (m *migrationRepository) GetByName(name string) (*models.Migration, error) {
	sd := &models.Migration{}
	if err := m.Conn.Where(&models.Migration{Name: name}).First(sd).Error; err != nil {
		return nil, m.wrapError(err)
	}

	return sd, nil
}

// Store creates the given migration and returns it's id.
func (m *migrationRepository) Store(sd *models.Migration) error {
	if err := m.Conn.Create(sd).Error; err != nil {
		return m.wrapError(err)
	}

	return nil
}

// DeleteByName deletes the given migration.
func (m *migrationRepository) DeleteByName(name string) error {
	if rowsAffected := m.Conn.Where(&models.Migration{Name: name}).
		Delete(&models.Migration{}).RowsAffected; rowsAffected == 0 {
		return m.wrapError(gorm.ErrRecordNotFound)
	}

	return nil
}

// wrapError wraps database error.
func (m *migrationRepository) wrapError(err error) error {
	return helpers.WrapError(err, tableName)
}
