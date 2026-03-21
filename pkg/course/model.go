package course

import (
	"errors"
	"fmt"
)

var (
	ErrValidation         = errors.New("validation error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrLoginExists        = errors.New("login already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrStudentNotFound    = errors.New("student not found")
)

const (
	defaultTokenTTLSeconds = 24 * 60 * 60
	bearerTokenType        = "Bearer"
	jwtAlgHS256            = "HS256"
	jwtTyp                 = "JWT"
)

type RegisterInput struct {
	Login     string
	Password  string
	Email     string
	FirstName string
	LastName  string
}

type LoginInput struct {
	Login    string
	Password string
}

type AuthToken struct {
	AccessToken string
	ExpiresIn   int
	TokenType   string
}

type AuthConfig struct {
	JWTSecret     string
	JWTTTLSeconds int
}

type ValidationError struct {
	Field  string
	Reason string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s %s", e.Field, e.Reason)
}

func (e ValidationError) Unwrap() error {
	return ErrValidation
}

type tokenHeader struct {
	Alg string
	Typ string
}

type tokenClaims struct {
	Sub   string
	Login string
	Exp   int64
	Iat   int64
}
