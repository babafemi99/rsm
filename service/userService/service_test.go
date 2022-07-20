package userService

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rsm/entity/userModel"
	"rsm/repository/userRepo"
	"testing"
)

var log = logrus.New()

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Persist(user *userModel.UserModel) (*userModel.UserAccessModel, error) {
	args := m.Called()
	results := args.Get(0)
	return results.(*userModel.UserAccessModel), args.Error(1)
}

func (m *MockRepository) Update(user *userModel.UserModel) (*userModel.UserModel, error) {
	args := m.Called()
	results := args.Get(0)
	return results.(*userModel.UserModel), args.Error(1)
}

func (m *MockRepository) Delete(id uuid.UUID) error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockRepository) FindById(id uuid.UUID) (*userModel.UserAccessModel, error) {
	args := m.Called(id)
	results := args.Get(0)
	return results.(*userModel.UserAccessModel), args.Error(1)
}

func (m *MockRepository) FindByEmail(email string) (*userModel.UserModel, error) {
	args := m.Called(email)
	results := args.Get(0)
	return results.(*userModel.UserModel), args.Error(1)
}

func TestGetById(t *testing.T) {
	id := uuid.New()

	mockRepo := new(MockRepository)

	userdata := userModel.UserAccessModel{
		Id:        id,
		FirstName: "bait",
		LastName:  "uus",
		Email:     "b@b.com",
	}

	mockRepo.On("FindById", id).Return(&userdata, nil)

	testSrv := NewUserService(log, mockRepo)
	user, err := testSrv.GetByUserId(id)

	assert.Nil(t, err)
	assert.Equal(t, id, user.Id)
}

func TestGetByWrongId(t *testing.T) {
	wrongUuid := uuid.New()

	mockRepo := new(MockRepository)

	mockRepo.On("FindById", wrongUuid).Return(&userModel.UserAccessModel{}, errors.New("no data found"))

	testSrv := NewUserService(log, mockRepo)
	user, err := testSrv.GetByUserId(wrongUuid)

	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "no data found", err.Error())
}

func Test_userService_Login(t *testing.T) {
	id := uuid.New()
	userdata := userModel.UserModel{
		Id:        id,
		FirstName: "bait",
		LastName:  "uus",
		Email:     "b@b.com",
		Password:  "$2a$14$2djvlayweuaxkot0fEbIsOOePfQ6Oer/IZSSb6qjSEp08gNSe8nnu",
	}

	reqCorrectCredentials := userModel.UserLoginRequest{
		Email:    "ooluwa27@gmail.com",
		Password: "secret1234",
	}
	reqInavalidEmail := userModel.UserLoginRequest{
		Email:    "ooluwa27com",
		Password: "secret1234",
	}

	mockRepo := new(MockRepository)

	mockRepo.On("FindByEmail", reqCorrectCredentials.Email).Return(&userdata, nil)
	mockRepo.On("FindByEmail", reqInavalidEmail.Email).Return(&userModel.UserModel{}, nil)

	type fields struct {
		log  *logrus.Logger
		repo userRepo.RepoInterface
	}
	type args struct {
		request userModel.UserLoginRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected *userModel.UserAccessModel
		err      error
	}{
		{
			name: "Correct credentials",
			fields: fields{
				log:  log,
				repo: mockRepo,
			},
			args: args{request: reqCorrectCredentials},
			expected: &userModel.UserAccessModel{
				Id:        id,
				FirstName: "bait",
				LastName:  "uus",
				Email:     "b@b.com",
			},
			err: nil,
		},
		{
			name: "Invalid Email",
			fields: fields{
				log:  log,
				repo: mockRepo,
			},
			args: args{
				request: reqInavalidEmail,
			},
			expected: nil,
			err:      fmt.Errorf("Key: 'UserLoginRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			u := NewUserService(tt.fields.log, tt.fields.repo)
			user, _ := u.Login(tt.args.request)
			assert.Equalf(t, tt.expected, user, "Login using : %v", tt.args.request)
		})
	}

}
