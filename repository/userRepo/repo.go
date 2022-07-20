package userRepo

import (
	"github.com/google/uuid"
	"rsm/entity/userModel"
)

type RepoInterface interface {
	Persist(user *userModel.UserModel) (*userModel.UserAccessModel, error)
	Update(user *userModel.UserModel) (*userModel.UserModel, error)
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) (*userModel.UserAccessModel, error)
	FindByEmail(email string) (*userModel.UserModel, error)
}
