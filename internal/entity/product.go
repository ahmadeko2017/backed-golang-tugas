package entity

import (
	"time"

	"github.com/microcosm-cc/bluemonday"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null;index" json:"name" binding:"required"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price" binding:"required"`
	Stock       int            `gorm:"not null" json:"stock" binding:"required"`
	CategoryID  uint           `gorm:"not null;index" json:"category_id" binding:"required"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.Name = bluemonday.UGCPolicy().Sanitize(p.Name)
	p.Description = bluemonday.UGCPolicy().Sanitize(p.Description)
	return
}

func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	p.Name = bluemonday.UGCPolicy().Sanitize(p.Name)
	p.Description = bluemonday.UGCPolicy().Sanitize(p.Description)
	return
}
