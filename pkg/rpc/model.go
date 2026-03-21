package rpc

type AuthToken struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

type RegisterRequest struct {
	Login     string `json:"login" validate:"required,min=3,max=64"`
	Password  string `json:"password" validate:"required,min=6,max=255"`
	Email     string `json:"email" validate:"required,email,max=255"`
	FirstName string `json:"firstName" validate:"required,max=255"`
	LastName  string `json:"lastName" validate:"required,max=255"`
}

type RegisterResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required,max=64"`
	Password string `json:"password" validate:"required,max=255"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

type StudentResponse struct {
	StudentID int    `json:"studentId"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type MeResponse struct {
	StudentID int    `json:"studentId"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
