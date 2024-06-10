package entity

import (
	"testing"
)

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		expected bool
	}{
		{
			name:     "Non-admin user",
			user:     &User{ID: 1, Login: "test", Password: "pass", Salt: nil, Admin: false},
			expected: false,
		},
		{
			name:     "Admin user",
			user:     &User{ID: 1, Login: "admin", Password: "pass", Salt: nil, Admin: true},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.user.IsAdmin(); result != tt.expected {
				t.Errorf("IsAdmin() = %v, want %v", result, tt.expected)
			}
		})
	}
}
