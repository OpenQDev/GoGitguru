// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Repo struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Organization string    `json:"organization"`
	Repository   string    `json:"repository"`
	Url          string    `json:"url"`
}
