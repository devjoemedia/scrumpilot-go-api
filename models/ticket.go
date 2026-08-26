package models

import (
	"time"

	"gorm.io/gorm"
)

type CreateTicketRequest struct {
	Title        string `json:"title" validate:"required,min=1,max=255"`
	Description  string `json:"description" validate:"omitempty,max=1000"`
	Status       string `json:"status" validate:"oneof=open in_progress closed"`
	Priority     string `json:"priority" validate:"oneof=low medium high"`
	AssigneeID   *uint  `json:"assignee_id" validate:"omitempty,gt=0"`
	AssignedByID *uint  `json:"assigned_by_id" validate:"omitempty,gt=0"`
}

type UpdateTicketRequest struct {
	Title        *string `json:"title" validate:"omitempty,min=1,max=255"`
	Description  *string `json:"description" validate:"omitempty,max=1000"`
	Status       *string `json:"status" validate:"omitempty,oneof=open in_progress closed"`
	Priority     *string `json:"priority" validate:"omitempty,oneof=low medium high"`
	AssigneeID   *uint   `json:"assignee_id" validate:"omitempty,gt=0"`
	AssignedByID *uint   `json:"assigned_by_id" validate:"omitempty,gt=0"`
}
type AssignTicketRequest struct {
	AssigneeID *uint `json:"assignee_id" validate:"required,gt=0"`
}

type Ticket struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Title       string `json:"title" validate:"required,min=5,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	Status      string `json:"status" validate:"oneof=open in_progress closed"`
	Priority    string `json:"priority" validate:"oneof=low medium high"`

	AssigneeID *uint `json:"assignee_id"`
	Assignee   *User `gorm:"foreignKey:AssigneeID" json:"assignee"`

	// Who created the ticket
	CreatedByID uint `json:"created_by_id"`
	CreatedBy   User `gorm:"foreignKey:CreatedByID"`

	// Who assigned the ticket
	AssignedByID *uint `json:"assigned_by_id"`
	AssignedBy   *User `gorm:"foreignKey:AssignedByID" json:"assigned_by"`

	Comments []Comment `gorm:"foreignKey:TicketID" json:"comments"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
