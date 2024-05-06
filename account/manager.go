package account

import (
	"context"
	"errors"
	"time"
	"time-tracker-backend/x/xjwt"

	"github.com/golang-jwt/jwt/v5"
)

type AccountManager struct {
	Persistence AccountPersistence
}

func NewManager(persistence AccountPersistence) *AccountManager {
	return &AccountManager{
		Persistence: persistence,
	}
}

func (p *AccountManager) Registration(ctx context.Context, name, email, password string) (string, string, error) {
	us, err := NewUser(name, email, "", password)
	if err != nil {
		return "", "", err
	}

	err = p.Persistence.Create(us)
	if err != nil {
		return "", "", err
	}

	return p.GenerateToken(ctx, us)
}

type JwTT struct {
	AccessToken  string
	RefreshToken string
}

func (p *AccountManager) Login(ctx context.Context, email string, password string) (string, string, error) {
	us, err := p.Persistence.Get(email)
	// validar erro de usuario nao existe
	if err != nil {
		return "", "", err
	}

	// TODO: adicionar criptografia
	if us.Password != password {
		return "", "", errors.New("password does not match")
	}

	return p.GenerateToken(ctx, us)
}

func (p *AccountManager) RefreshToken(refreshToken string) (string, string, error) {
	err := xjwt.VerifyToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	us, err := p.Persistence.GetUserByRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	token, refreshToken, err := xjwt.TokenAndRefreshToken(jwt.MapClaims{
		"id":    us.ID,
		"name":  us.Name,
		"email": us.Email,
		"kind":  "session",
	}, time.Hour*2, time.Hour*24)
	if err != nil {
		return "", "", err
	}

	us.Token = token
	us.RefreshToken = refreshToken

	err = p.Persistence.Save(us)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (p *AccountManager) GenerateToken(ctx context.Context, us *User) (tokenString string, refreshToken string, err error) {
	token, refreshToken, err := xjwt.TokenAndRefreshToken(jwt.MapClaims{
		"id":    us.ID,
		"name":  us.Name,
		"email": us.Email,
		"kind":  "session",
	}, time.Hour*2, time.Hour*24)
	if err != nil {
		return "", "", err
	}

	us.Token = token
	us.RefreshToken = refreshToken

	err = p.Persistence.Save(us)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (p *AccountManager) GetUserByToken(ctx context.Context, token string) (*User, error) {
	return p.Persistence.GetUserByToken(token)
}
