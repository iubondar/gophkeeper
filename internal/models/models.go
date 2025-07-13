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

// SecretData представляет данные секрета
type SecretData struct {
	Name string
	Type string
	Data string
}

// TextSecretData представляет данные текстового секрета
type TextSecretData struct {
	Name string
	Text string
}

// LoginPasswordData представляет данные логина и пароля
type LoginPasswordData struct {
	Name     string
	Login    string
	Password string
	URL      string
}

// CardData представляет данные банковской карты
type CardData struct {
	Name   string
	Number string
	Holder string
	Expiry string
	CVV    string
}

// FileData представляет данные файла
type FileData struct {
	Name     string
	FilePath string
}
