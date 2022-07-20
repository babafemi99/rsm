package userModel

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserModel struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8,alphanum"`
	CreatedAt time.Time `json:"-"`
}

type UserAccessModel struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Email     string    `json:"email" validate:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,alphanum"`
}

func (u *UserModel) ValidateInput() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserLoginRequest) ValidateInput() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserModel) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *UserModel) ComparePasswords(s string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(s))
}
