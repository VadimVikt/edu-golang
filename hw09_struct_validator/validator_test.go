package hw09structvalidator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"  //nolint
	"github.com/stretchr/testify/require" //nolint
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin, stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "User_Invalid_Email_Fails",
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "Ivan",
				Age:    25,
				Email:  "invalid email",
				Role:   "admin",
				Phones: []string{"79991234567"},
			},
			expectedErr: ErrInvalidFormat,
		},
		{
			name: "App_Invalid_Version_Length_Fails",
			in: App{
				Version: "0.0.1-1",
			},
			expectedErr: ErrInvalidLength,
		},
		{
			name: "Token_No_Validation_Tags_Passes",
			in: Token{
				Header:    []byte("Authorization"),
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			name: "Response_Valid_Code_404_Passes",
			in: Response{
				Code: 405,
				Body: "Not Found",
			},
			expectedErr: ErrValueOutOfRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)

			if tt.expectedErr == nil {
				assert.Equal(t, err.Error(), "no errors", tt.name)
			} else {
				require.Error(t, err, "Expected error for %s", tt.name)
				assert.ErrorIs(t, err, tt.expectedErr, "Error mismatch for %s", tt.name)
			}
		})
	}
}

func TestNoStructType(t *testing.T) {
	user := 1
	err := Validate(user)

	t.Logf("Actual error type: %T", err)
	t.Logf("Actual error: %+v", err)
	require.Error(t, err, "Expected error for %s")
	assert.ErrorIs(t, err, ErrExpectedStruct, "Error mismatch for %s")
}
