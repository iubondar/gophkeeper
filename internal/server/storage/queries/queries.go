package queries

const (
	InsertUser = `
		INSERT INTO users (id, login, password_hash, salt)
		VALUES ($1, $2, $3, $4)
	`

	GetUserSalt = `
		SELECT salt FROM users WHERE login = $1
	`
)
