// Package migration describes migration interfaces.
package migration

import (
	"github.com/sulimak-co/gorm-migrator/internal/models"
)

// Repository interface represents the migration's repository contract.
type Repository interface {
	GetByName(name string) (*models.Migration, error)
	Store(migration *models.Migration) error
	DeleteByName(name string) error
}
