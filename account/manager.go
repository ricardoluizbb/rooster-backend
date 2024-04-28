package account

import (
	"errors"
	"fmt"
	"net/url"
	"time"
	"time-tracker-backend/config"
	"time-tracker-backend/x/xjwt"

	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret = []byte("secretKeyHere") //TODO: em variavel de ambiente

type AccountManager struct {
	Persistence AccountPersistence
}

func NewManager(persistence AccountPersistence) *AccountManager {
	return &AccountManager{
		Persistence: persistence,
	}
}

func (p *AccountManager) CreateUser(name, email string) (string, error) {
	us, err := NewUser(name, email, "")
	if err != nil {
		return "", err
	}

	err = p.Persistence.Create(us)
	if err != nil {
		return "", err
	}

	magicLink, err := GenerateMagicLink(us.ID, us.Email)
	if err != nil {
		return "", err
	}

	return magicLink, nil
}

type JwTT struct {
	AccessToken  string
	RefreshToken string
}

func (p *AccountManager) Login(email string) (string, error) {
	us, err := p.Persistence.Get(email)
	// validar erro de usuario nao existe
	if err != nil {
		return "", err
	}

	magicLink, err := GenerateMagicLink(us.ID, us.Email)
	if err != nil {
		return "", err
	}

	return magicLink, nil
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

func (p *AccountManager) MagicLink(magicToken string) (tokenString string, refreshToken string, err error) {
	err = xjwt.VerifyToken(magicToken)
	if err != nil {
		return
	}

	token, err := jwt.Parse(magicToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		accountEmail := fmt.Sprintf("%v", claims["email"])

		us, err := p.Persistence.Get(accountEmail)
		if err != nil {
			return tokenString, refreshToken, err
		}

		return xjwt.TokenAndRefreshToken(jwt.MapClaims{
			"id":    us.ID,
			"name":  us.Name,
			"email": us.Email,
		}, time.Hour*2, time.Hour*24)

	} else {
		err = errors.New("cannot parse token")
		return
	}
}

func GenerateMagicLink(id string, email string) (string, error) {
	magicToken, err := xjwt.Token(jwt.MapClaims{
		"id":    id,
		"email": email,
	}, time.Minute*10, "magic")
	if err != nil {
		return "", err
	}

	u, _ := url.Parse(config.MagicLinkUrl())

	query := url.Values{}
	query.Set("magicToken", magicToken)

	u.RawQuery = query.Encode()

	return u.String(), nil
}
