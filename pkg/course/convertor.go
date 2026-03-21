package course

import (
	"strconv"

	"courses/pkg/db"
)

func newRegisterStudent(in RegisterInput, passwordHash string) *db.Student {
	return &db.Student{
		Login:        in.Login,
		PasswordHash: passwordHash,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Email:        in.Email,
		StatusID:     db.StatusEnabled,
	}
}

func newAuthToken(token string, expiresIn int) AuthToken {
	return AuthToken{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   bearerTokenType,
	}
}

func newTokenHeader() tokenHeader {
	return tokenHeader{
		Alg: jwtAlgHS256,
		Typ: jwtTyp,
	}
}

func newTokenClaims(studentID int, login string, iat, exp int64) tokenClaims {
	return tokenClaims{
		Sub:   strconv.Itoa(studentID),
		Login: login,
		Exp:   exp,
		Iat:   iat,
	}
}
