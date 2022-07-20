package userModel

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserModel_WrongValidationInput(t *testing.T) {
	type fields struct {
		Id        uuid.UUID
		FirstName string
		LastName  string
		Email     string
		Password  string
		createdAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{

		{
			name: "no firstName",
			fields: fields{
				Id:        uuid.New(),
				FirstName: "",
				LastName:  "bayo",
				Email:     "bayo@bayo.com",
				Password:  "babanlabayo",
				createdAt: time.Now(),
			},
		}, {
			name: "no lastname",
			fields: fields{
				Id:        uuid.New(),
				FirstName: "bayo",
				LastName:  "",
				Email:     "bayo@bayo.com",
				Password:  "babanlabayo",
				createdAt: time.Now(),
			},
		}, {
			name: "no email",
			fields: fields{
				Id:        uuid.New(),
				FirstName: "ade",
				LastName:  "bayo",
				Email:     "",
				Password:  "babanlabayo",
				createdAt: time.Now(),
			},
		}, {
			name: "bad email format",
			fields: fields{
				Id:        uuid.New(),
				FirstName: "ade",
				LastName:  "bayo",
				Email:     "bayo.com",
				Password:  "babanlabayo",
				createdAt: time.Now(),
			},
		}, {
			name: "no password",
			fields: fields{
				Id:        uuid.New(),
				FirstName: "ade",
				LastName:  "bayo",
				Email:     "bayo@b.com",
				Password:  "",
				createdAt: time.Now(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserModel{
				Id:        tt.fields.Id,
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
				CreatedAt: tt.fields.createdAt,
			}

			err := u.ValidateInput()
			assert.NotNil(t, err)
		})
	}
}
