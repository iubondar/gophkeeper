package models

type RegisterIn struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	Salt         string `json:"salt"`
}

type LoginIn struct {
	Login string `json:"login"`
}

type LoginOut struct {
	Salt string `json:"salt"`
}

type AuthenticateIn struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
}

type AuthenticateOut struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthenticateResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	Success      bool
}

type UserCredentials struct {
	Login    string
	Password string
}
