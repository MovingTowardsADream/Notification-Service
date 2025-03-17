package hasher_test

import (
	"testing"

	"Notification_Service/pkg/hasher"
)

func TestSHA1Hasher_Hash(t *testing.T) {
	testCases := []struct {
		name     string
		salt     string
		password string
		expected string
	}{
		{
			name:     "simple password with salt",
			salt:     "somesalt",
			password: "password123",
			expected: "736f6d6573616c74cbfdac6008f9cab4083784cbd1874f76618d2a97",
		},
		{
			name:     "empty password with salt",
			salt:     "somesalt",
			password: "",
			expected: "736f6d6573616c74da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		{
			name:     "password with empty salt",
			salt:     "",
			password: "password123",
			expected: "cbfdac6008f9cab4083784cbd1874f76618d2a97",
		},
		{
			name:     "empty password and empty salt",
			salt:     "",
			password: "",
			expected: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := hasher.NewSHA1Hasher(tc.salt)
			result := h.Hash(tc.password)

			if result != tc.expected {
				t.Errorf("Test case '%s' failed: expected '%s', got '%s'", tc.name, tc.expected, result)
			}
		})
	}
}
