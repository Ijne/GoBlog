package tools

import (
	"os"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func TestValidateEmail(t *testing.T) {
	cases := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "correct email",
			email:    "example@smth.com",
			expected: true,
		},
		{
			name:     "incorrect email",
			email:    "example.smth.com",
			expected: false,
		},
		{
			name:     "empty",
			email:    "",
			expected: false,
		},
		{
			name:     "space",
			email:    " ",
			expected: false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := ValidateEmail(tc.email)
			if got != tc.expected {
				t.Errorf("Email(%v) = %v; want %v", tc.email, got, tc.expected)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		t.Errorf("Error %v", err)
	}
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	cases := []struct {
		name        string
		token       string
		wantErr     bool
		errContains string
	}{
		{
			name: "correct signingmethod",
			token: func() string {
				tk, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{}).SignedString(jwtSecret)
				return tk
			}(),
			wantErr:     false,
			errContains: "",
		},
		{
			name: "incorrect signingmethod 384",
			token: func() string {
				tk, _ := jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{"name": "Ivan"}).SignedString(jwtSecret)
				return tk
			}(),
			wantErr:     true,
			errContains: "token contains an invalid number of segments",
		},
		{
			name: "incorrect signingmethod 512",
			token: func() string {
				tk, _ := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{"name": "Ivan"}).SignedString(jwtSecret)
				return tk
			}(),
			wantErr:     true,
			errContains: "token contains an invalid number of segments",
		},
		{
			name:        "empty token",
			token:       "",
			wantErr:     true,
			errContains: "invalid number of segments",
		},
		{
			name:        "invalid token format",
			token:       "invalid.token.format",
			wantErr:     true,
			errContains: "invalid character '\\u008a' looking for beginning of value",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ValidateToken(tc.token)

			hasErr := err != nil
			if hasErr != tc.wantErr {
				t.Errorf("ValidateToken() error existence = %v, want %v", hasErr, tc.wantErr)
				return
			}

			if tc.wantErr && err != nil {
				if tc.errContains != "" && !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("ValidateToken() error = %q, want contains %q", err.Error(), tc.errContains)
				}
			}

			if !tc.wantErr && err != nil {
				t.Errorf("ValidateToken() unexpected error = %v", err)
			}
		})
	}
}
