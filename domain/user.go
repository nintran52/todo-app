package domain

import (
	"errors"
	"strings"
	"time"
	"todo-app/pkg/clients"

	"github.com/google/uuid"
)

const EntityName = "User"

type UserRole int

const (
	RoleUser UserRole = 1 << iota
	RoleAdmin
)

func (role UserRole) String() string {
	switch role {
	case RoleAdmin:
		return "admin"
	default:
		return "user"
	}
}

type User struct {
	ID        uuid.UUID
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Phone     string         `json:"phone"`
	Role      UserRole       `json:"role"`
	Salt      string         `json:"-"`
	Status    clients.Status `json:"status"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

type UserCreate struct {
	ID        uuid.UUID
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Role      UserRole `json:"-"`
	Salt      string   `json:"-"`
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}

func (ic *UserCreate) Validate() error {
	var validationErrors []string

	if ic.Email == "" {
		validationErrors = append(validationErrors, "email can not be null")
	}
	if ic.Password == "" {
		validationErrors = append(validationErrors, "password can not be null")
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}

	return nil
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

var (
	ErrEmailOrPasswordInvalid = clients.NewCustomError(
		errors.New("email or password invalid"),
		"email or password invalid",
		"ErrUsernameOrPasswordInvalid",
	)

	ErrEmailExisted = clients.NewCustomError(
		errors.New("email has already existed"),
		"email has already existed",
		"ErrEmailExisted",
	)
)
