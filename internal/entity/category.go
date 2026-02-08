package entity

import (
	"time"

	"github.com/microcosm-cc/bluemonday"

	"gorm.io/gorm"
)

// Category represents the category model
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name" binding:"required"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.Name = bluemonday.UGCPolicy().Sanitize(c.Name)
	c.Description = bluemonday.UGCPolicy().Sanitize(c.Description)
	return
}

func (c *Category) BeforeUpdate(tx *gorm.DB) (err error) {
	c.Name = bluemonday.UGCPolicy().Sanitize(c.Name)
	c.Description = bluemonday.UGCPolicy().Sanitize(c.Description)
	return
}
