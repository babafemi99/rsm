package passwordUtils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	ComparePasswords(plain, s string) error
	HashPassword(password string) (string, error)
}

type passwordSev struct {
	log *logrus.Logger
}

func (p passwordSev) ComparePasswords(plain, s string) error {
	byteHash := []byte(s)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plain))
	if err != nil {
		p.log.Errorf("Error comparing passwords: %v", err)
		return err
	}
	fmt.Println("error isnt nil oh ")
	return nil
}

func (p passwordSev) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		p.log.Errorf("Error generating password: %v", err)
		return "", err
	}

	return string(bytes), nil
}
