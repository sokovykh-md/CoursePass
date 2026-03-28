package rpc

import (
	"context"
	"net/mail"

	"courses/pkg/coursepass"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type AuthService struct {
	zenrpc.Service
	embedlog.Logger

	authManager *coursepass.AuthManager
}

func NewAuthService(dbc db.DB, logger embedlog.Logger, authCfg coursepass.AuthConfig) *AuthService {
	return &AuthService{
		authManager: coursepass.NewAuthManager(dbc, logger, authCfg),
		Logger:      logger,
	}
}

func (as *AuthService) Register(ctx context.Context, login, password, email, firstName, lastName string) (*Token, error) {
	if err := validateRegisterRequest(login, password, email, firstName, lastName); err != nil {
		as.Logger.Error(ctx, "auth register invalid params", "err", err)
		return nil, err
	}

	token, err := as.authManager.Register(ctx, login, password, email, firstName, lastName)
	if err != nil {
		as.Logger.Error(ctx, "auth register failed", "err", err)
		return nil, mapRPCError(err)
	}

	return newToken(token), nil
}

func (as *AuthService) Login(ctx context.Context, login, password string) (*Token, error) {
	if err := validateLoginRequest(login, password); err != nil {
		as.Logger.Error(ctx, "auth login invalid params", "err", err)
		return nil, err
	}

	token, err := as.authManager.Login(ctx, login, password)
	if err != nil {
		as.Logger.Error(ctx, "auth login failed", "err", err)
		return nil, mapRPCError(err)
	}

	return newToken(token), nil
}

func validateRegisterRequest(login, password, email, firstName, lastName string) error {
	if login == "" {
		return invalidParamsError("login", "is required")
	}
	if len([]rune(login)) > 255 {
		return invalidParamsError("login", "max length is 255")
	}

	if password == "" {
		return invalidParamsError("password", "is required")
	}
	if len([]rune(password)) < 6 {
		return invalidParamsError("password", "min length is 6")
	}
	if len([]rune(password)) > 255 {
		return invalidParamsError("password", "max length is 255")
	}

	if email == "" {
		return invalidParamsError("email", "is required")
	}
	if len([]rune(email)) > 255 {
		return invalidParamsError("email", "max length is 255")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return invalidParamsError("email", "invalid format")
	}

	if firstName == "" {
		return invalidParamsError("firstName", "is required")
	}
	if len([]rune(firstName)) > 255 {
		return invalidParamsError("firstName", "max length is 255")
	}

	if lastName == "" {
		return invalidParamsError("lastName", "is required")
	}
	if len([]rune(lastName)) > 255 {
		return invalidParamsError("lastName", "max length is 255")
	}

	return nil
}

func validateLoginRequest(login, password string) error {
	if login == "" {
		return invalidParamsError("login", "is required")
	}
	if len([]rune(login)) > 255 {
		return invalidParamsError("login", "max length is 255")
	}

	if password == "" {
		return invalidParamsError("password", "is required")
	}
	if len([]rune(password)) > 255 {
		return invalidParamsError("password", "max length is 255")
	}

	return nil
}
