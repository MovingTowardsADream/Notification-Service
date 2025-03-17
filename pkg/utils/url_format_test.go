package utils_test

import (
	"testing"

	"Notification_Service/pkg/utils"
)

func TestFormatAddress(t *testing.T) {
	testCases := []struct {
		name     string
		host     string
		port     int
		expected string
	}{
		{
			name:     "valid host and port",
			host:     "localhost",
			port:     8080,
			expected: "localhost:8080",
		},
		{
			name:     "empty host",
			host:     "",
			port:     8080,
			expected: ":8080",
		},
		{
			name:     "zero port",
			host:     "localhost",
			port:     0,
			expected: "localhost:0",
		},
		{
			name:     "large port",
			host:     "localhost",
			port:     65535,
			expected: "localhost:65535",
		},
		{
			name:     "ipv6 address",
			host:     "::1",
			port:     8080,
			expected: "::1:8080",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.FormatAddress(tc.host, tc.port)
			if result != tc.expected {
				t.Errorf("Test case '%s' failed: expected '%s', got '%s'", tc.name, tc.expected, result)
			}
		})
	}
}
