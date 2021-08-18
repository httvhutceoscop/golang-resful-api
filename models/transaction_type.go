package models

import "github.com/gofrs/uuid"

type TransactionType struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}
