package psqlRepo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"rsm/entity/userModel"
	"rsm/repository/userRepo"
)

type psql struct {
	log  *logrus.Logger
	conn *pgx.Conn
}

func (p *psql) Update(user *userModel.UserModel) (*userModel.UserModel, error) {
	persistStmt := fmt.Sprintf(
		"UPDATE User SET id = %v, firstname =%v, lastname=%v, email=%v",
		user.Id, user.FirstName, user.LastName, user.Email)

	_, err := p.conn.Exec(context.Background(), persistStmt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *psql) Persist(user *userModel.UserModel) (*userModel.UserAccessModel, error) {
	persistStmt := fmt.Sprintf(
		"INSERT INTO User (id, firstname, lastname, email, created_at) VALUES(%v, %v, %v, %v, %v", user.Id,
		user.FirstName, user.LastName, user.Email, user.CreatedAt)

	_, err := p.conn.Exec(context.Background(), persistStmt)
	if err != nil {
		return nil, err
	}
	userAccess := userModel.UserAccessModel{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	return &userAccess, nil

}

func (p *psql) Delete(id uuid.UUID) error {
	deleteStmt := fmt.Sprintf("DELETE FROM User WHERE id=%v", id)
	_, err := p.conn.Exec(context.Background(), deleteStmt)
	if err != nil {
		return err
	}
	return nil
}

func (p *psql) FindById(id uuid.UUID) (*userModel.UserAccessModel, error) {
	var userAccess userModel.UserAccessModel
	findByIdStmt := fmt.Sprintf("SELECT id, firstname, lastname, email FROM User WHERE id = %v", id)
	queryRes, err := p.conn.Query(context.Background(), findByIdStmt)
	if err != nil {
		p.log.Errorf("Error Finding By Id: %v", err)
		return nil, err
	}
	err = queryRes.Scan(userAccess.Id, userAccess.FirstName, userAccess.LastName, userAccess.Email)
	if err != nil {
		p.log.Errorf("Error Scanning into Access Model: %v", err)
		return nil, err
	}
	return &userAccess, nil
}

func (p *psql) FindByEmail(email string) (*userModel.UserModel, error) {
	var userAccess userModel.UserModel
	findByEmailStmt := fmt.Sprintf("SELECT id, firstname, lastname, email, password FROM User WHERE email = %v", email)
	queryRes, err := p.conn.Query(context.Background(), findByEmailStmt)
	if err != nil {
		p.log.Errorf("Error Finding By Email: %v", err)
		return nil, err
	}
	err = queryRes.Scan(userAccess.Id, userAccess.FirstName, userAccess.LastName, userAccess.Email, userAccess.Password)
	if err != nil {
		p.log.Errorf("Error Scanning into Access Model: %v", err)
		return nil, err
	}
	return &userAccess, nil
}

func NewPsqlService(conn *pgx.Conn, log *logrus.Logger) userRepo.RepoInterface {
	return &psql{conn: conn, log: log}
}
