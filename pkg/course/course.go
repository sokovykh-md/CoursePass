package course

import (
	"context"
	"fmt"

	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"golang.org/x/crypto/bcrypt"
)

type CourseManager struct {
	dbc    db.DB
	tlRepo db.CoursesRepo
	auth   AuthConfig
	embedlog.Logger
}

func NewCourseManager(dbc db.DB, logger embedlog.Logger, authCfg AuthConfig) *CourseManager {
	return &CourseManager{
		dbc:    dbc,
		tlRepo: db.NewCoursesRepo(dbc),
		auth:   authCfg,
		Logger: logger,
	}
}

func (cm *CourseManager) Register(ctx context.Context, in RegisterInput) (AuthToken, error) {
	if student, err := cm.tlRepo.OneStudent(ctx, &db.StudentSearch{Login: &in.Login}); err != nil {
		return AuthToken{}, err
	} else if student != nil {
		return AuthToken{}, ErrLoginExists
	}

	if student, err := cm.tlRepo.OneStudent(ctx, &db.StudentSearch{Email: &in.Email}); err != nil {
		return AuthToken{}, err
	} else if student != nil {
		return AuthToken{}, ErrEmailExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed generate hash password: %w", err)
	}

	student, err := cm.tlRepo.AddStudent(ctx, newRegisterStudent(in, string(passwordHash)))
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create student: %w", err)
	}

	token, expiresIn, err := generateJWT(cm.auth, student.ID, student.Login)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create JWT: %w", err)
	}

	return newAuthToken(token, expiresIn), nil
}

func (cm *CourseManager) Login(ctx context.Context, in LoginInput) (AuthToken, error) {
	student, err := cm.tlRepo.OneStudent(ctx, &db.StudentSearch{
		Login: &in.Login,
	})
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed get student: %w", err)
	}
	if student == nil {
		return AuthToken{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(student.PasswordHash), []byte(in.Password)); err != nil {
		return AuthToken{}, ErrInvalidCredentials
	}

	token, expiresIn, err := generateJWT(cm.auth, student.ID, student.Login)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create JWT: %w", err)
	}

	return newAuthToken(token, expiresIn), nil
}

func (cm *CourseManager) Me(ctx context.Context, studentID int) (*db.Student, error) {
	student, err := cm.tlRepo.OneStudent(ctx, &db.StudentSearch{
		ID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get student: %w", err)
	}
	if student == nil {
		return nil, ErrStudentNotFound
	}

	return student, nil
}
