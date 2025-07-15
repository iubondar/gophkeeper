package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mustMarshal маршалит объект в JSON для тестов
func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	assert.NoError(t, err)
	return data
}
