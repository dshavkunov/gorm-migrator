package models

import "time"

// Migration describes seed model.
type Migration struct {
	ID        int32     `json:"id" gorm:"primary_key;" faker:"-"`
	Name      string    `json:"name" gorm:"size:255;not null;unique" validate:"required"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}
