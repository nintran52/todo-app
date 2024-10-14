package clients

import "github.com/google/uuid"

type TokenPayload struct {
	UID   uuid.UUID `json:"user_id"`
	URole string    `json:"role"`
}

func (p TokenPayload) UserID() uuid.UUID {
	return p.UID
}

func (p TokenPayload) Role() string {
	return p.URole
}

type Requester interface {
	GetUserID() uuid.UUID
	GetEmail() string
	GetRole() string
}
