package models

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestJSONError_Error(t *testing.T) {
	err := &JSONError{Message: "test error", Code: 400}
	if err.Error() != "test error" {
		t.Errorf("expected 'test error', got '%s'", err.Error())
	}
}

func TestNewJSONError(t *testing.T) {
	err := NewJSONError("msg", 404)
	if err.Message != "msg" || err.Code != 404 {
		t.Errorf("unexpected JSONError: %+v", err)
	}
}

func TestJSONError_Encode(t *testing.T) {
	rec := httptest.NewRecorder()
	err := &JSONError{Message: "encode error", Code: 418}
	err.Encode(rec)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != 418 {
		t.Errorf("expected status 418, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected content-type application/json, got %s", ct)
	}

	var got JSONError
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Message != "encode error" || got.Code != 418 {
		t.Errorf("unexpected JSONError: %+v", got)
	}
}

func TestEncodeError(t *testing.T) {
	rec := httptest.NewRecorder()
	EncodeError(rec, "encode error", 401)
	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != 401 {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
	var got JSONError
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Message != "encode error" || got.Code != 401 {
		t.Errorf("unexpected JSONError: %+v", got)
	}
}

func TestParseJSONError(t *testing.T) {
	orig := &JSONError{Message: "parse me", Code: 500}
	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	parsed, err := ParseJSONError(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed.Message != orig.Message || parsed.Code != orig.Code {
		t.Errorf("expected %+v, got %+v", orig, parsed)
	}
}

func TestParseJSONError_Invalid(t *testing.T) {
	b := []byte(`not a json`)
	_, err := ParseJSONError(b)
	if err == nil {
		t.Error("expected error for invalid json, got nil")
	}
}
