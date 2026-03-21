package rpc

import (
	"courses/pkg/course"
	"courses/pkg/db"
)

func newRegisterInput(req RegisterRequest) course.RegisterInput {
	return course.RegisterInput{
		Login:     req.Login,
		Password:  req.Password,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
}

func newLoginInput(req LoginRequest) course.LoginInput {
	return course.LoginInput{
		Login:    req.Login,
		Password: req.Password,
	}
}

func newRegisterResponse(token course.AuthToken) RegisterResponse {
	return RegisterResponse{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		TokenType:   token.TokenType,
	}
}

func newLoginResponse(token course.AuthToken) LoginResponse {
	return LoginResponse{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		TokenType:   token.TokenType,
	}
}

func newMeResponse(student *db.Student) MeResponse {
	return MeResponse{
		StudentID: student.ID,
		Login:     student.Login,
		Email:     student.Email,
		FirstName: student.FirstName,
		LastName:  student.LastName,
	}
}
