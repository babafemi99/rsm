package userService

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"rsm/crypto/passwordUtils"
	"rsm/entity/userModel"
	"rsm/repository/userRepo"
	"time"
)

type ServiceInterface interface {
	Login(request userModel.UserLoginRequest) (*userModel.UserAccessModel, error)
	SignUp(model *userModel.UserModel) (*userModel.UserAccessModel, error)
	GetByEmail(email string) (*userModel.UserAccessModel, error)
	GetByUserId(id uuid.UUID) (*userModel.UserAccessModel, error)
	DeleteUser()
	GetAllUsers()
}

func (u *userService) Login(request userModel.UserLoginRequest) (*userModel.UserAccessModel, error) {
	err := request.ValidateInput()
	if err != nil {
		u.log.Errorf("Validation Error: %v", err)
		return nil, fmt.Errorf("something went wrong while validation")
	}

	accessUser, findingErr := u.repo.FindByEmail(request.Email)
	if findingErr != nil {
		return nil, findingErr
	}

	err = u.crypto.ComparePasswords(request.Password, accessUser.Password)
	if err != nil {
		u.log.Errorf("Password Validation Error: %v", err)
		return nil, fmt.Errorf("invalid password")
	}

	userAccess := userModel.UserAccessModel{
		Id:        accessUser.Id,
		FirstName: accessUser.FirstName,
		LastName:  accessUser.LastName,
		Email:     accessUser.Email,
	}
	return &userAccess, nil
}

func (u *userService) SignUp(model *userModel.UserModel) (*userModel.UserAccessModel, error) {
	err := model.ValidateInput()
	if err != nil {
		u.log.Errorf("Validation Error: %v", err)
		return nil, fmt.Errorf("something went wrong while validation")
	}
	_, err = u.repo.FindByEmail(model.Email)
	if errors.Is(err, nil) {
		return nil, fmt.Errorf("user already exists")
	}

	password, cryptErr := u.crypto.HashPassword(model.Password)
	if cryptErr != nil {
		return nil, cryptErr
	}
	model.Password = password
	model.CreatedAt = time.Now()
	return u.repo.Persist(model)
}

func (u *userService) GetByEmail(email string) (*userModel.UserAccessModel, error) {
	accessUser, err := u.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	userAccess := userModel.UserAccessModel{
		Id:        accessUser.Id,
		FirstName: accessUser.FirstName,
		LastName:  accessUser.LastName,
		Email:     accessUser.Email,
	}
	return &userAccess, nil
}

func (u *userService) GetByUserId(id uuid.UUID) (*userModel.UserAccessModel, error) {
	u.log.Info("Inside User Get Service")
	accessUser, err := u.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	userAccess := userModel.UserAccessModel{
		Id:        accessUser.Id,
		FirstName: accessUser.FirstName,
		LastName:  accessUser.LastName,
		Email:     accessUser.Email,
	}
	return &userAccess, nil
}

func (u *userService) DeleteUser() {
	//TODO implement me
	panic("implement me")
}

func (u *userService) GetAllUsers() {
	//TODO implement me
	panic("implement me")
}

type userService struct {
	log    *logrus.Logger
	repo   userRepo.RepoInterface
	crypto passwordUtils.PasswordService
}

func NewUserService(log *logrus.Logger, repo userRepo.RepoInterface, c passwordUtils.PasswordService) ServiceInterface {
	return &userService{log: log, repo: repo, crypto: c}
}
