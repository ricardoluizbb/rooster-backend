package models

type LoginRequest struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	SetCookie bool   `json:"setCookie,omitempty"`
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken,omitempty"`
}

type Tokens struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RegistrationRequest struct {
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	ConfirmPassword string `json:"confirmPassword,omitempty"`
	SetCookie       bool   `json:"setCookie,omitempty"`
}
