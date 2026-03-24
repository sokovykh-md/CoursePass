package coursepass

import (
	"errors"
	"testing"

	"courses/pkg/db"
	dbtest "courses/pkg/db/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthManager_Register_Success(t *testing.T) {
	// Arrange
	dbo, logger := dbtest.Setup(t)
	authCfg := AuthConfig{
		JWTSecret:     "test-secret",
		JWTTTLSeconds: 3600,
	}
	manager := NewAuthManager(dbo, logger, authCfg)
	repo := db.NewCoursesRepo(dbo)

	login := "student_" + dbtest.NextStringID()
	email := "student_" + dbtest.NextStringID() + "@mail.test"

	// Act
	token, err := manager.Register(t.Context(), login, "password123", email, "John", "Doe")

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.Equal(t, authCfg.JWTTTLSeconds, token.ExpiresIn)
	assert.Equal(t, "Bearer", token.TokenType)

	studentID, err := ValidateJWT(authCfg, token.AccessToken)
	require.NoError(t, err)
	assert.Positive(t, studentID)

	student, err := repo.OneStudent(t.Context(), &db.StudentSearch{ID: &studentID})
	require.NoError(t, err)
	require.NotNil(t, student)
	assert.Equal(t, login, student.Login)
	assert.Equal(t, email, student.Email)
	assert.Equal(t, "John", student.FirstName)
	assert.Equal(t, "Doe", student.LastName)
}

func TestAuthManager_Register_DuplicateLogin(t *testing.T) {
	// Arrange
	dbo, logger := dbtest.Setup(t)
	manager := NewAuthManager(dbo, logger, AuthConfig{JWTSecret: "test-secret"})

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	existingLogin := "student_" + dbtest.NextStringID()
	_, cleanup := dbtest.Student(t, dbo.DB, &db.Student{
		Login:        existingLogin,
		PasswordHash: string(passwordHash),
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "student_" + dbtest.NextStringID() + "@mail.test",
		StatusID:     1,
	})
	defer cleanup()

	// Act
	_, err = manager.Register(
		t.Context(),
		existingLogin,
		"password123",
		"student_"+dbtest.NextStringID()+"@mail.test",
		"John",
		"Doe",
	)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrLoginExists))
}

func TestAuthManager_Register_DuplicateEmail(t *testing.T) {
	// Arrange
	dbo, logger := dbtest.Setup(t)
	manager := NewAuthManager(dbo, logger, AuthConfig{JWTSecret: "test-secret"})

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	existingEmail := "student_" + dbtest.NextStringID() + "@mail.test"
	_, cleanup := dbtest.Student(t, dbo.DB, &db.Student{
		Login:        "student_" + dbtest.NextStringID(),
		PasswordHash: string(passwordHash),
		FirstName:    "John",
		LastName:     "Doe",
		Email:        existingEmail,
		StatusID:     1,
	})
	defer cleanup()

	// Act
	_, err = manager.Register(
		t.Context(),
		"student_"+dbtest.NextStringID(),
		"password123",
		existingEmail,
		"John",
		"Doe",
	)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmailExists))
}

func TestAuthManager_Login_Success(t *testing.T) {
	// Arrange
	dbo, logger := dbtest.Setup(t)
	authCfg := AuthConfig{JWTSecret: "test-secret"}
	manager := NewAuthManager(dbo, logger, authCfg)

	rawPassword := "password123"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	require.NoError(t, err)

	student, cleanup := dbtest.Student(t, dbo.DB, &db.Student{
		Login:        "student_" + dbtest.NextStringID(),
		PasswordHash: string(passwordHash),
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "student_" + dbtest.NextStringID() + "@mail.test",
		StatusID:     1,
	})
	defer cleanup()

	// Act
	token, err := manager.Login(t.Context(), student.Login, rawPassword)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.Equal(t, "Bearer", token.TokenType)

	studentID, err := ValidateJWT(authCfg, token.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, student.ID, studentID)
}

func TestAuthManager_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	dbo, logger := dbtest.Setup(t)
	manager := NewAuthManager(dbo, logger, AuthConfig{JWTSecret: "test-secret"})

	rawPassword := "password123"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	require.NoError(t, err)

	student, cleanup := dbtest.Student(t, dbo.DB, &db.Student{
		Login:        "student_" + dbtest.NextStringID(),
		PasswordHash: string(passwordHash),
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "student_" + dbtest.NextStringID() + "@mail.test",
		StatusID:     1,
	})
	defer cleanup()

	// Act
	_, err = manager.Login(t.Context(), student.Login, "wrong-password")

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidCredentials))
}
