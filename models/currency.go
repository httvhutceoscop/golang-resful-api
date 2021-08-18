package models

import "github.com/gofrs/uuid"

type Currency struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}
