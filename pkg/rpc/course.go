package rpc

import (
	"context"
	"courses/pkg/course"
	"courses/pkg/db"
	"net/mail"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type CourseService struct {
	zenrpc.Service
	embedlog.Logger

	courseManager *course.CourseManager
}

func NewCourseService(dbc db.DB, logger embedlog.Logger, authCfg course.AuthConfig) *CourseService {
	return &CourseService{
		courseManager: course.NewCourseManager(dbc, logger, authCfg),
		Logger:        logger,
	}
}

func (cs *CourseService) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		cs.Logger.Error(ctx, "auth register invalid params", "err", err)
		return RegisterResponse{}, err
	}

	token, err := cs.courseManager.Register(ctx, newRegisterInput(req))
	if err != nil {
		cs.Logger.Error(ctx, "auth register failed", "err", err)
		return RegisterResponse{}, mapRPCError(err)
	}

	return newRegisterResponse(token), nil
}

func (cs *CourseService) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	if err := validateLoginRequest(req); err != nil {
		cs.Logger.Error(ctx, "auth login invalid params", "err", err)
		return LoginResponse{}, err
	}

	token, err := cs.courseManager.Login(ctx, newLoginInput(req))
	if err != nil {
		cs.Logger.Error(ctx, "auth login failed", "err", err)
		return LoginResponse{}, mapRPCError(err)
	}

	return newLoginResponse(token), nil
}

func (cs *CourseService) Me(ctx context.Context) (MeResponse, error) {
	studentID, ok := StudentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		cs.Logger.Error(ctx, "auth me failed: no studentID in context")
		return MeResponse{}, mapRPCError(course.ErrInvalidToken)
	}

	student, err := cs.courseManager.Me(ctx, studentID)
	if err != nil {
		cs.Logger.Error(ctx, "auth me failed", "err", err)
		return MeResponse{}, mapRPCError(err)
	}

	return newMeResponse(student), nil
}

func validateRegisterRequest(req RegisterRequest) error {
	if req.Login == "" {
		return invalidParamsError("login", "is required")
	}
	if len([]rune(req.Login)) > 255 {
		return invalidParamsError("login", "max length is 255")
	}

	if req.Password == "" {
		return invalidParamsError("password", "is required")
	}
	if len([]rune(req.Password)) < 6 {
		return invalidParamsError("password", "min length is 6")
	}
	if len([]rune(req.Password)) > 255 {
		return invalidParamsError("password", "max length is 255")
	}

	if req.Email == "" {
		return invalidParamsError("email", "is required")
	}
	if len([]rune(req.Email)) > 255 {
		return invalidParamsError("email", "max length is 255")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return invalidParamsError("email", "invalid format")
	}

	if req.FirstName == "" {
		return invalidParamsError("firstName", "is required")
	}
	if len([]rune(req.FirstName)) > 255 {
		return invalidParamsError("firstName", "max length is 255")
	}

	if req.LastName == "" {
		return invalidParamsError("lastName", "is required")
	}
	if len([]rune(req.LastName)) > 255 {
		return invalidParamsError("lastName", "max length is 255")
	}

	return nil
}

func validateLoginRequest(req LoginRequest) error {
	if req.Login == "" {
		return invalidParamsError("login", "is required")
	}
	if len([]rune(req.Login)) > 255 {
		return invalidParamsError("login", "max length is 255")
	}

	if req.Password == "" {
		return invalidParamsError("password", "is required")
	}
	if len([]rune(req.Password)) > 255 {
		return invalidParamsError("password", "max length is 255")
	}

	return nil
}
