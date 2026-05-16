package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	UserId    uuid.UUID `json:"user_id"`
}
