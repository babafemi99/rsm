package userService

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
		return nil, err
	}
	accessUser, findingErr := u.repo.FindByEmail(request.Email)
	if findingErr != nil {
		return nil, findingErr
	}
	err = accessUser.ComparePasswords(request.Password)
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

func (u *userService) SignUp(model *userModel.UserModel) (*userModel.UserAccessModel, error) {
	err := model.ValidateInput()
	if err != nil {
		return nil, err
	}
	oldUser, _ := u.repo.FindByEmail(model.Email)
	if oldUser != nil {
		return nil, errors.New("user Already Exits")
	}
	err = model.HashPassword()
	if err != nil {
		return nil, err
	}
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
	log  *logrus.Logger
	repo userRepo.RepoInterface
}

func NewUserService(log *logrus.Logger, repo userRepo.RepoInterface) ServiceInterface {
	return &userService{log: log, repo: repo}
}
