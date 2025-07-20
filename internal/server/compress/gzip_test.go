package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzipCompression(t *testing.T) {
	requestBody := `
		<html><body><h1>Hello world!</h1></body></html>
	`

	successBody := `{
		"result": "https://127.0.0.1/abcdef11"
	}`

	withoutGzip := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentType, "application/json")
		io.WriteString(w, successBody)
	})
	handler := WithGzipCompression(withoutGzip)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set(contentEncoding, "gzip")
		r.Header.Set(acceptEncoding, "")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.JSONEq(t, successBody, string(b))
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		r := httptest.NewRequest("POST", srv.URL, buf)
		r.RequestURI = ""
		r.Header.Set(contentType, "text/html")
		r.Header.Set(acceptEncoding, "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		assert.JSONEq(t, successBody, string(b))
	})
}

func TestGzipWriter_NoDuplicateWriteHeader(t *testing.T) {
	// Создаем тестовый ResponseWriter
	recorder := httptest.NewRecorder()

	// Создаем gzipWriter
	gw := newGzipWriter(recorder)

	// Устанавливаем заголовки
	gw.Header().Set("Content-Type", "application/json")

	// Вызываем WriteHeader первый раз
	gw.WriteHeader(http.StatusOK)

	// Пытаемся вызвать WriteHeader второй раз - это не должно вызвать ошибку
	gw.WriteHeader(http.StatusInternalServerError)

	// Проверяем, что статус остался первым установленным
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	// Проверяем, что заголовок Content-Encoding установлен
	if recorder.Header().Get("Content-Encoding") != "gzip" {
		t.Errorf("Expected Content-Encoding: gzip, got %s", recorder.Header().Get("Content-Encoding"))
	}

	// Закрываем writer
	gw.Close()
}

func TestGzipWriter_WriteBeforeWriteHeader(t *testing.T) {
	// Создаем тестовый ResponseWriter
	recorder := httptest.NewRecorder()

	// Создаем gzipWriter
	gw := newGzipWriter(recorder)

	// Устанавливаем заголовки
	gw.Header().Set("Content-Type", "application/json")

	// Записываем данные до WriteHeader
	_, err := gw.Write([]byte("test data"))
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}

	// Проверяем, что заголовок Content-Encoding установлен после Write
	if recorder.Header().Get("Content-Encoding") != "gzip" {
		t.Errorf("Expected Content-Encoding: gzip, got %s", recorder.Header().Get("Content-Encoding"))
	}

	// Закрываем writer
	gw.Close()
}
