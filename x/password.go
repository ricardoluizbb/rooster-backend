package x

import (
	"errors"
	"unicode"
)

func IsValidPassword(rawPassword string) error {
	if len(rawPassword) < 6 || len(rawPassword) > 64 {
		return errors.New("password length is invalid. required >=6 and <=64")
	}

	var hasUpperCase, hasLowerCase, hasNumber, hasSymbol bool
	for _, r := range rawPassword {
		if unicode.IsUpper(r) {
			hasUpperCase = true
		}
		if unicode.IsLower(r) {
			hasLowerCase = true
		}
		if unicode.IsNumber(r) {
			hasNumber = true
		}
		if unicode.IsSymbol(r) || unicode.IsPunct(r) {
			hasSymbol = true
		}
	}

	isValid := hasUpperCase && hasLowerCase && hasNumber && hasSymbol
	if !isValid {
		// TODO: melhorar esse erro
		return errors.New("password is invalid")
	}

	return nil
}
