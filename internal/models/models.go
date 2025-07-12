package models

type RegisterIn struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	Salt         string `json:"salt"`
}

type LoginIn struct {
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
}

type UserCredentials struct {
	Login    string
	Password string
}
