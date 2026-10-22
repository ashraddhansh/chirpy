package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "Valid bearer token",
			header:    "Bearer abc123",
			wantToken: "abc123",
			wantErr:   false,
		},
		{
			name:      "Bearer scheme is lowercase",
			header:    "bearer abc123",
			wantToken: "abc123",
			wantErr:   false,
		},
		{
			name:      "Extra whitespace",
			header:    "  Bearer   abc123  ",
			wantToken: "abc123",
			wantErr:   false,
		},
		{
			name:    "Missing authorization header",
			header:  "",
			wantErr: true,
		},
		{
			name:    "Missing token",
			header:  "Bearer",
			wantErr: true,
		},
		{
			name:    "Invalid authorization scheme",
			header:  "Basic abc123",
			wantErr: true,
		},
		{
			name:    "Missing bearer scheme",
			header:  "abc123",
			wantErr: true,
		},
		{
			name:    "Too many fields",
			header:  "Bearer abc123 extra",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(http.Header)
			if tt.header != "" {
				headers.Set("Authorization", tt.header)
			}

			token, err := GetBearerToken(headers)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && token != tt.wantToken {
				t.Errorf("GetBearerToken() = %v, want %v", token, tt.wantToken)
			}
		})
	}
}
