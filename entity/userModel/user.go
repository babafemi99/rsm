package userModel

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type UserModel struct {
	Id        uuid.UUID `json:"id" validator:"required"`
	FirstName string    `json:"firstName" validator:"required"`
	LastName  string    `json:"lastName" validator:"required"`
	Email     string    `json:"email" validator:"required"`
	Password  string    `json:"password" validator:"required,min=8,alphanum"`
	createdAt time.Time `json:"created_at"`
}

type UserAccessModel struct {
	Id        uuid.UUID `json:"id" validator:"required"`
	FirstName string    `json:"firstName" validator:"required"`
	LastName  string    `json:"lastName" validator:"required"`
	Email     string    `json:"email" validator:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserModel) ValidateInput() error {
	validate := validator.New()
	return validate.Struct(u)
}
