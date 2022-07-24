package userService

import (
	"bou.ke/monkey"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rsm/crypto/passwordUtils"
	"rsm/entity/userModel"
	"rsm/repository/userRepo"
	"testing"
	"time"
)

var log = logrus.New()

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Persist(user *userModel.UserModel) (*userModel.UserAccessModel, error) {
	args := m.Called(user)
	results := args.Get(0)
	return results.(*userModel.UserAccessModel), args.Error(1)
}

func (m *MockRepository) Update(user *userModel.UserModel) (*userModel.UserModel, error) {
	args := m.Called(user)
	results := args.Get(0)
	return results.(*userModel.UserModel), args.Error(1)
}

func (m *MockRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
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

type mockPasswordUtils struct {
	mock.Mock
}

func (m *mockPasswordUtils) ComparePasswords(plain, s string) error {
	args := m.Called(plain, s)
	return args.Error(0)
}

func (m *mockPasswordUtils) HashPassword(password string) (string, error) {
	args := m.Called(password)
	results := args.String(0)
	return results, args.Error(1)
}

func TestGetById(t *testing.T) {
	id := uuid.New()

	mockRepo := new(MockRepository)
	mockPass := new(mockPasswordUtils)

	userdata := userModel.UserAccessModel{
		Id:        id,
		FirstName: "bait",
		LastName:  "uus",
		Email:     "b@b.com",
	}

	mockRepo.On("FindById", id).Return(&userdata, nil)

	testSrv := NewUserService(log, mockRepo, mockPass)
	user, err := testSrv.GetByUserId(id)

	assert.Nil(t, err)
	assert.Equal(t, id, user.Id)
}

func TestGetByWrongId(t *testing.T) {
	wrongUuid := uuid.New()

	mockRepo := new(MockRepository)
	mockPass := new(mockPasswordUtils)

	mockRepo.On("FindById", wrongUuid).Return(&userModel.UserAccessModel{}, errors.New("no data found"))

	testSrv := NewUserService(log, mockRepo, mockPass)
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
	reqInavalidPassword := userModel.UserLoginRequest{
		Email:    "ooluwa27@gmail.com",
		Password: "ozy",
	}

	rewWrongCredentials := userModel.UserLoginRequest{
		Email:    "ooluwa27@gmails.com",
		Password: "secret12345",
	}

	mockRepo := new(MockRepository)
	mockPass := new(mockPasswordUtils)

	mockRepo.On("FindByEmail", reqCorrectCredentials.Email).Return(&userdata, nil)
	mockRepo.On("FindByEmail", reqInavalidEmail.Email).Return(&userModel.UserModel{}, fmt.Errorf("error Finding User"))
	mockRepo.On("FindByEmail", reqInavalidPassword.Email).Return(&userModel.UserModel{},
		fmt.Errorf("error Finding User"))
	mockRepo.On("FindByEmail", rewWrongCredentials.Email).Return(&userModel.UserModel{},
		fmt.Errorf("error Finding User"))

	mockPass.On("ComparePasswords", reqCorrectCredentials.Password,
		"$2a$14$2djvlayweuaxkot0fEbIsOOePfQ6Oer/IZSSb6qjSEp08gNSe8nnu").Return(nil)

	mockPass.On("ComparePasswords", rewWrongCredentials.Password,
		"$2a$14$2djvlayweuaxkot0fEbIsOOePfQ6Oer/IZSSb6qjSEp08gNSe8nnu").Return(fmt.Errorf("wrong credentials"))

	mockPass.On("ComparePasswords", reqInavalidPassword.Password,
		"$2a$14$2djvlayweuaxkot0fEbIsOOePfQ6Oer/IZSSb6qjSEp08gNSe8nnu").Return(fmt.Errorf("error invalid password"))

	mockPass.On("ComparePasswords", reqInavalidEmail.Password,
		"$2a$14$2djvlayweuaxkot0fEbIsOOePfQ6Oer/IZSSb6qjSEp08gNSe8nnu").Return(fmt.Errorf("error invalid password"))

	type fields struct {
		log    *logrus.Logger
		repo   userRepo.RepoInterface
		crypto passwordUtils.PasswordService
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
				log:    log,
				repo:   mockRepo,
				crypto: mockPass,
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
				log:    log,
				repo:   mockRepo,
				crypto: mockPass,
			},
			args: args{
				request: reqInavalidEmail,
			},
			expected: nil,
		}, {
			name: "Invalid password",
			fields: fields{
				log:  log,
				repo: mockRepo,
			},
			args: args{
				request: reqInavalidPassword,
			},
			expected: nil,
		}, {
			name: "Wrong Email and Password Combination",
			fields: fields{
				log:  log,
				repo: mockRepo,
			},
			args: args{
				request: rewWrongCredentials,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			u := NewUserService(tt.fields.log, tt.fields.repo, tt.fields.crypto)
			user, _ := u.Login(tt.args.request)
			assert.Equalf(t, tt.expected, user, "Login using : %v", tt.args.request)
		})
	}

}

func Test_userService_GetByEmail(t *testing.T) {
	id := uuid.New()
	email := "b@b.com"
	wrongemail := "wrong@email.com"

	mockRepo := new(MockRepository)

	userdata := userModel.UserModel{
		Id:        id,
		FirstName: "bait",
		LastName:  "uus",
		Email:     "b@b.com",
		Password:  "$2a$14$2djvlayweuaxkot0fEbIsOOePfQ6Oer/IZSSb6qjSEp08gNSe8nnu",
	}

	mockRepo.On("FindByEmail", email).Return(&userdata, nil)
	mockRepo.On("FindByEmail", wrongemail).Return(&userModel.UserModel{}, fmt.Errorf("error finding by email"))

	type fields struct {
		log  *logrus.Logger
		repo userRepo.RepoInterface
	}
	type args struct {
		email string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *userModel.UserAccessModel
	}{
		{
			name: "correct credentials",
			fields: fields{
				log:  log,
				repo: mockRepo,
			},
			args: args{
				email: email,
			},
			want: &userModel.UserAccessModel{
				Id:        id,
				FirstName: "bait",
				LastName:  "uus",
				Email:     "b@b.com",
			},
		}, {
			name: "wrong credentials",
			fields: fields{
				log:  log,
				repo: mockRepo,
			},
			args: args{
				email: wrongemail,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userService{
				log:  tt.fields.log,
				repo: tt.fields.repo,
			}
			user, _ := u.GetByEmail(tt.args.email)

			assert.Equalf(t, tt.want, user, "GetByEmail(%v)", tt.args.email)
		})
	}
}

func Test_userService_SignUp(t *testing.T) {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	})
	fmt.Println(time.Now())
	id := uuid.New()
	model1 := userModel.UserModel{
		Id:        id,
		FirstName: "Ab",
		LastName:  "Cd",
		Email:     "abdc@abcd.com",
		Password:  "secret12345",
		CreatedAt: time.Now(),
	}
	persistmodel := userModel.UserModel{
		Id:        id,
		FirstName: "Ab",
		LastName:  "Cd",
		Email:     "abdc@abcd.com",
		Password:  "$2a$14$2djvlayweuaxkot0fEbIsOOePfQ6Oer/IZSSb6qjSEp08gNSe8nnu",
		CreatedAt: time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	model1access := userModel.UserAccessModel{
		Id:        id,
		FirstName: "Ab",
		LastName:  "Cd",
		Email:     "abdc@abcd.com",
	}

	model2 := userModel.UserModel{
		Id:        id,
		FirstName: "ef",
		LastName:  "gh",
		Email:     "efgh@efgj.com",
		Password:  "secret12345",
	}

	mockRepo := new(MockRepository)
	mockPass := new(mockPasswordUtils)

	mockRepo.On("FindByEmail", model1.Email).Return(&userModel.UserModel{}, fmt.Errorf("error finding by email"))
	mockRepo.On("FindByEmail", model2.Email).Return(&model2, nil)

	mockRepo.On("Persist", &persistmodel).Return(&model1access, nil)

	mockPass.On("HashPassword", model1.Password).Return(
		"$2a$14$2djvlayweuaxkot0fEbIsOOePfQ6Oer/IZSSb6qjSEp08gNSe8nnu", nil)

	type fields struct {
		log    *logrus.Logger
		repo   userRepo.RepoInterface
		crypto passwordUtils.PasswordService
	}
	type args struct {
		model *userModel.UserModel
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *userModel.UserAccessModel
		err    error
	}{
		{
			name: "correct credentials",
			fields: fields{
				log:    log,
				repo:   mockRepo,
				crypto: mockPass,
			},
			args: args{
				model: &model1,
			},
			want: &model1access,
			err:  nil,
		},
		{
			name: "existing user ",
			fields: fields{
				log:    log,
				repo:   mockRepo,
				crypto: mockPass,
			},
			args: args{
				model: &model2,
			},
			want: nil,
			err:  fmt.Errorf("user already exists"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userService{
				log:    tt.fields.log,
				repo:   tt.fields.repo,
				crypto: tt.fields.crypto,
			}
			got, _ := u.SignUp(tt.args.model)

			assert.Same(t, tt.want, got, "SignUp(%v)", tt.args.model)
		})
	}
}

func Test_userService_DeleteUser(t *testing.T) {
	id := uuid.New()

	mockRepo := new(MockRepository)
	mockPass := new(mockPasswordUtils)
	mockRepo.On("Delete", id).Return(nil)

	type fields struct {
		log    *logrus.Logger
		repo   userRepo.RepoInterface
		crypto passwordUtils.PasswordService
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "delete user",
			fields: fields{
				log:    log,
				repo:   mockRepo,
				crypto: mockPass,
			},
			args: args{
				id: id,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userService{
				log:    tt.fields.log,
				repo:   tt.fields.repo,
				crypto: tt.fields.crypto,
			}
			got := u.DeleteUser(tt.args.id)
			assert.Equalf(t, tt.wantErr, got, "DeleteUser(%v)", tt.args.id)
		})
	}
}
