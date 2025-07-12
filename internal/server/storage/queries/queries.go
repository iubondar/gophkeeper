package queries

const (
	InsertUser = `
		INSERT INTO users (id, login, password_hash, salt)
		VALUES ($1, $2, $3, $4)
	`
)
