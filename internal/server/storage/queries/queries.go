package queries

const (
	InsertUser = `
		INSERT INTO users (id, login, password_hash, salt)
		VALUES ($1, $2, $3, $4)
	`

	GetUserSalt = `
		SELECT salt FROM users WHERE login = $1
	`

	GetUserByLoginAndPassword = `
		SELECT id FROM users WHERE login = $1 AND password_hash = $2
	`

	InsertRecord = `
		INSERT INTO records (id, user_id, label, type, metadata, encrypted_data, file_key, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	GetRecordByLabel = `
		SELECT id, label, type, metadata, encrypted_data, file_key, version, created_at, updated_at
		FROM records 
		WHERE label = $1 AND user_id = $2
	`

	DeleteRecordByLabel = `
		DELETE FROM records WHERE label = $1 AND user_id = $2
	`
)
