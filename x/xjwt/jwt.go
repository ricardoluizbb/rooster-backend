package xjwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var hmacSampleSecret = []byte("secretKeyHere") //TODO: em variavel de ambiente

func TokenAndRefreshToken(claims jwt.MapClaims, tokenExpiresInHour, refreshTokenExpiresInHour time.Duration) (string, string, error) {

	claims["expiresAt"] = jwt.NewNumericDate(time.Now().Add(tokenExpiresInHour))
	claims["kind"] = "session"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", "", err
	}

	claims["kind"] = "refresh"
	claims["expiresAt"] = jwt.NewNumericDate(time.Now().Add(refreshTokenExpiresInHour))
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshTokenString, err := refreshToken.SignedString(hmacSampleSecret)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func Token(claims jwt.MapClaims, tokenExpiresInHour time.Duration, kind string) (string, error) {
	claims["expiresAt"] = jwt.NewNumericDate(time.Now().Add(tokenExpiresInHour))
	claims["kind"] = kind
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(hmacSampleSecret)
}
