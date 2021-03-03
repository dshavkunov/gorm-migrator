package migrations

import (
	"gorm.io/gorm"
	"time"
)

//createUserTableUp creates user table with belongs to organization relation
func createUserTableUp(migrator gorm.Migrator) error {
	// organization used for relation
	type organization struct {
		ID string `json:"id" gorm:"type:uuid;primary_key;" faker:"uuid_hyphenated"`
	}

	type user struct {
		ID             string    `json:"id" gorm:"type:uuid;primary_key" faker:"uuid_hyphenated"`
		Email          string    `json:"email" gorm:"size:100;not null;unique" validate:"required,email"`
		Role           int       `json:"role" gorm:"default:0"`
		OrganizationID string    `json:"organizationID" gorm:"type:uuid;not null" validate:"required" faker:"uuid_hyphenated"`
		Password       string    `json:"password" gorm:"size:100;not null;"  validate:"required"`
		Username       string    `json:"username" gorm:"size:255;not null;unique" validate:"required"`
		CreatedAt      time.Time `json:"createdAt" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
		UpdatedAt      time.Time `json:"updatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
		// related entities
		Organization organization `json:"-" validate:"-" faker:"-"`
	}

	return migrator.CreateTable(&user{})
}

//createUserTableDown deletes user table
func createUserTableDown(migrator gorm.Migrator) error {
	type user struct{}
	return migrator.DropTable(&user{})
}
