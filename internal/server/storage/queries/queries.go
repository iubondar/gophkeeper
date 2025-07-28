// Package queries содержит SQL-запросы для работы с базой данных PostgreSQL.
// Пакет содержит константы с SQL-запросами для операций с пользователями
// и секретами, включая вставку, выборку, обновление и удаление записей.
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
		INSERT INTO records (id, user_id, label, type, metadata, encrypted_data, file_key, file_name, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	GetRecordByLabel = `
		SELECT id, label, type, metadata, encrypted_data, file_key, file_name, version, created_at, updated_at
		FROM records 
		WHERE label = $1 AND user_id = $2
	`

	DeleteRecordByLabel = `
		DELETE FROM records WHERE label = $1 AND user_id = $2
	`

	// Обновление записи с проверкой версии (OCC)
	UpdateRecordByLabelAndVersion = `
	    UPDATE records
	    SET type = $1, metadata = $2, encrypted_data = $3, file_key = $4, file_name = $5, version = version + 1, updated_at = $6
	    WHERE label = $7 AND user_id = $8 AND version = $9
	`
)
