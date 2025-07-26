package validators

import (
	"testing"

	"gophkeeper/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestValidateUserCredentials(t *testing.T) {
	tests := []struct {
		name        string
		credentials *models.UserCredentials
		wantOk      bool
		wantErr     bool
	}{
		{
			name: "valid credentials",
			credentials: &models.UserCredentials{
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
			credentials: &models.UserCredentials{
				Login:    "",
				Password: "testpass",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "empty password",
			credentials: &models.UserCredentials{
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
		data    *models.TextSecretData
		wantOk  bool
		wantErr bool
	}{
		{
			name: "valid data",
			data: &models.TextSecretData{
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
			data: &models.TextSecretData{
				Name: "",
				Text: "secret text",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "empty text",
			data: &models.TextSecretData{
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
			secretName: "testsecret",
			wantOk:     true,
			wantErr:    false,
		},
		{
			name:       "empty name",
			secretName: "",
			wantOk:     false,
			wantErr:    true,
		},
		{
			name:       "invalid name with space",
			secretName: "test secret",
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
		data    *models.CardData
		wantOk  bool
		wantErr bool
	}{
		{
			name: "valid card",
			data: &models.CardData{
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
			data: &models.CardData{
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
			data: &models.CardData{
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
			data: &models.CardData{
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
			data: &models.CardData{
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
			data: &models.CardData{
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
			data: &models.CardData{
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
			data: &models.CardData{
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

func TestValidateCardData_Luhn(t *testing.T) {
	validCard := &models.CardData{
		Name:     "Test Card",
		Number:   "4111111111111111", // Valid Luhn number
		Holder:   "John Doe",
		Expiry:   "12/25",
		CVV:      "123",
		Metadata: "Test metadata",
	}

	ok, err := ValidateCardData(validCard)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestValidateTextSecretData_WithMetadata(t *testing.T) {
	validData := &models.TextSecretData{
		Name:     "Test Secret",
		Text:     "Secret text content",
		Metadata: "Test metadata",
	}

	ok, err := ValidateTextSecretData(validData)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestValidateLoginPasswordData_WithMetadata(t *testing.T) {
	validData := &models.LoginPasswordData{
		Name:     "Test Login",
		Login:    "user@example.com",
		Password: "password123",
		URL:      "https://example.com",
		Metadata: "Test metadata",
	}

	ok, err := ValidateLoginPasswordData(validData)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestValidateLoginPasswordData_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		data    *models.LoginPasswordData
		wantOk  bool
		wantErr bool
	}{
		{
			name: "empty name",
			data: &models.LoginPasswordData{
				Name:     "",
				Login:    "user",
				Password: "pass",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "empty login",
			data: &models.LoginPasswordData{
				Name:     "loginname",
				Login:    "",
				Password: "pass",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "empty password",
			data: &models.LoginPasswordData{
				Name:     "loginname",
				Login:    "user",
				Password: "",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "short name",
			data: &models.LoginPasswordData{
				Name:     "ab",
				Login:    "user",
				Password: "pass",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "short login",
			data: &models.LoginPasswordData{
				Name:     "loginname",
				Login:    "ab",
				Password: "pass",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "short password",
			data: &models.LoginPasswordData{
				Name:     "loginname",
				Login:    "user",
				Password: "ab",
			},
			wantOk:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateLoginPasswordData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLoginPasswordData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok != tt.wantOk {
				t.Errorf("ValidateLoginPasswordData() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}

func TestValidateFileData_WithMetadata(t *testing.T) {
	validData := &models.FileData{
		Name:     "Test File",
		FilePath: "/path/to/file.txt",
		Metadata: "Test metadata",
	}

	ok, err := ValidateFileData(validData)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestValidateFileData_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		data    *models.FileData
		wantOk  bool
		wantErr bool
	}{
		{
			name: "empty name",
			data: &models.FileData{
				Name:     "",
				FilePath: "/path/to/file.txt",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "short name",
			data: &models.FileData{
				Name:     "ab",
				FilePath: "/path/to/file.txt",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "empty file path",
			data: &models.FileData{
				Name:     "filename",
				FilePath: "",
			},
			wantOk:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := ValidateFileData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok != tt.wantOk {
				t.Errorf("ValidateFileData() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}
