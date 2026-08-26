package models

import (
	"time"

	"gorm.io/gorm"
)

type CreateCommentRequest struct {
	Comment  string `json:"comment" validate:"required,max=1000"`
	TicketID *uint  `json:"ticket_id" validate:"required,gt=0"`
}

type UpdateCommentRequest struct {
	Comment  string `json:"comment"`
	TicketID *uint  `json:"ticket_id" validate:"required,gt=0"`
}

type Comment struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Comment string `json:"comment" validate:"required,max=1000"`

	TicketID uint    `gorm:"column:ticket_id" json:"ticket_id"`
	Ticket   *Ticket `gorm:"foreignKey:TicketID" json:"ticket"`
	UserID   uint    `json:"user_id"`
	User     *User   `gorm:"foreignKey:UserID" json:"user"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
