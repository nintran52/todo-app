package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"-"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status" gorm:"column:status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (Item) TableName() string { return "items" }

type ItemCreation struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (ItemCreation) TableName() string { return Item{}.TableName() }

func (ic *ItemCreation) Validate() error {
	var validationErrors []string

	if ic.Title == "" {
		validationErrors = append(validationErrors, "title can not be null")
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}

	return nil
}

type ItemUpdate struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Status      *Status   `json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ItemUpdate) TableName() string { return Item{}.TableName() }
