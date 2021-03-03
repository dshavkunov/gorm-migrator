package migrations

import (
	"gorm.io/gorm"
	"time"
)

//createOrganizationTableUp creates organization table
func createOrganizationTableUp(migrator gorm.Migrator) error {
	type organization struct {
		ID          string    `json:"id" gorm:"type:uuid;primary_key;" faker:"uuid_hyphenated"`
		Name        string    `json:"name" gorm:"size:255;not null;unique" validate:"required"`
		Description string    `json:"description" gorm:"size:10000"`
		CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
		UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
	}

	return migrator.CreateTable(&organization{})
}

//createOrganizationTableUp deletes organization table
func createOrganizationTableDown(migrator gorm.Migrator) error {
	type organization struct{}
	return migrator.DropTable(&organization{})
}
