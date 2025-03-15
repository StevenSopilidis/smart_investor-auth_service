package domain

import "github.com/google/uuid"

type User struct {
	Id            uuid.UUID `gorm:primaryKey`
	Email         string
	EmaiLVerified bool
	Password      string
}
