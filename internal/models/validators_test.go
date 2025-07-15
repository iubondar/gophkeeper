package models

import (
	"testing"
)

func TestValidateUserCredentials(t *testing.T) {
	tests := []struct {
		name        string
		credentials *UserCredentials
		wantOk      bool
		wantErr     bool
	}{
		{
			name: "valid credentials",
			credentials: &UserCredentials{
				Login:    "testuser",
				Password: "testpass",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:        "nil credentials",
			credentials: nil,
			wantOk:      false,
			wantErr:     true,
		},
		{
			name: "empty login",
			credentials: &UserCredentials{
				Login:    "",
				Password: "testpass",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "empty password",
			credentials: &UserCredentials{
				Login:    "testuser",
				Password: "",
			},
			wantOk:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateUserCredentials(tt.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok != tt.wantOk {
				t.Errorf("ValidateUserCredentials() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}

func TestValidateTextSecretData(t *testing.T) {
	tests := []struct {
		name    string
		data    *TextSecretData
		wantOk  bool
		wantErr bool
	}{
		{
			name: "valid data",
			data: &TextSecretData{
				Name: "test secret",
				Text: "secret text",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:    "nil data",
			data:    nil,
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "empty name",
			data: &TextSecretData{
				Name: "",
				Text: "secret text",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "empty text",
			data: &TextSecretData{
				Name: "test secret",
				Text: "",
			},
			wantOk:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateTextSecretData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTextSecretData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok != tt.wantOk {
				t.Errorf("ValidateTextSecretData() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}

func TestValidateSecretName(t *testing.T) {
	tests := []struct {
		name       string
		secretName string
		wantOk     bool
		wantErr    bool
	}{
		{
			name:       "valid name",
			secretName: "test secret",
			wantOk:     true,
			wantErr:    false,
		},
		{
			name:       "empty name",
			secretName: "",
			wantOk:     false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateSecretName(tt.secretName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSecretName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok != tt.wantOk {
				t.Errorf("ValidateSecretName() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}

func TestValidateCardData(t *testing.T) {
	tests := []struct {
		name    string
		data    *CardData
		wantOk  bool
		wantErr bool
	}{
		{
			name: "valid card",
			data: &CardData{
				Name:   "test card",
				Number: "4003600000000014",
				Holder: "John Doe",
				Expiry: "12/25",
				CVV:    "123",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:    "nil data",
			data:    nil,
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "empty number",
			data: &CardData{
				Name:   "test card",
				Number: "",
				Holder: "John Doe",
				Expiry: "12/25",
				CVV:    "123",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "number not 16 digits",
			data: &CardData{
				Name:   "test card",
				Number: "12345678",
				Holder: "John Doe",
				Expiry: "12/25",
				CVV:    "123",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "number not only digits",
			data: &CardData{
				Name:   "test card",
				Number: "1234abcd12345678",
				Holder: "John Doe",
				Expiry: "12/25",
				CVV:    "123",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "invalid expiry format",
			data: &CardData{
				Name:   "test card",
				Number: "4003600000000014",
				Holder: "John Doe",
				Expiry: "1225",
				CVV:    "123",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "invalid expiry month",
			data: &CardData{
				Name:   "test card",
				Number: "4003600000000014",
				Holder: "John Doe",
				Expiry: "13/25",
				CVV:    "123",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "CVV not only digits",
			data: &CardData{
				Name:   "test card",
				Number: "4003600000000014",
				Holder: "John Doe",
				Expiry: "12/25",
				CVV:    "12a",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "CVV wrong length",
			data: &CardData{
				Name:   "test card",
				Number: "4003600000000014",
				Holder: "John Doe",
				Expiry: "12/25",
				CVV:    "12",
			},
			wantOk:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateCardData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCardData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok != tt.wantOk {
				t.Errorf("ValidateCardData() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}

func TestValidateLuhn(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid card number", "4003600000000014", true},
		{"Invalid card number", "4111111111111112", false},
		{"Empty string", "", false},
		{"Non-digit characters", "4111-1111-1111-1111", false},
		{"Single digit", "5", false},
		{"Single digit positive", "0", true},
		{"Two digits", "12", false},
		{"Two digits valid", "18", true},
		{"Long valid number", "4532015112830366", true},
		{"Long invalid number", "4532015112830367", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateLuhn(tt.input); got != tt.expected {
				t.Errorf("ValidateLuhn(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"user@example.com", true},
		{"user.name@domain.co", true},
		{"user@sub.domain.com", true},
		{"user@domain", false},
		{"userdomain.com", false},
		{"@domain.com", false},
		{"user@.com", false},
		{"user@domain.", false},
		{"user@domain.c", true}, // базовая проверка, не проверяет длину tld
		{"user@domain.corporate", true},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := isValidEmail(tt.email); got != tt.valid {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, got, tt.valid)
			}
		})
	}
}

func TestValidateCardData_Luhn(t *testing.T) {
	data := &CardData{
		Name:   "test card",
		Number: "4003600000000014", // валидный по Luhn
		Holder: "John Doe",
		Expiry: "12/25",
		CVV:    "123",
	}
	ok, err := ValidateCardData(data)
	if !ok || err != nil {
		t.Errorf("ValidateCardData() with valid Luhn failed: %v", err)
	}

	data.Number = "4111111111111112" // невалидный по Luhn
	ok, err = ValidateCardData(data)
	if ok || err == nil {
		t.Errorf("ValidateCardData() with invalid Luhn should fail")
	}
}
