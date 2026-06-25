package emailaddr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple valid email",
			input:    "user@gmail.com",
			expected: "user@gmail.com",
		},
		{
			name:     "uppercase domain gets lowercased",
			input:    "user@GMAIL.COM",
			expected: "user@gmail.com",
		},
		{
			name:     "uppercase local part gets lowercased",
			input:    "User@gmail.com",
			expected: "user@gmail.com",
		},
		{
			name:     "mixed case everywhere",
			input:    "UsEr@GmAil.CoM",
			expected: "user@gmail.com",
		},
		{
			name:     "leading and trailing whitespace trimmed",
			input:    "  user@gmail.com  ",
			expected: "user@gmail.com",
		},
		{
			name:     "local part with dot and hyphen",
			input:    "user.name-label@hotmail.com",
			expected: "user.name-label@hotmail.com",
		},
		{
			name:     "numeric local part after letter",
			input:    "a123@gmail.com",
			expected: "a123@gmail.com",
		},
		{
			name:     "subdomain in domain",
			input:    "user@mail.gmail.com",
			expected: "user@mail.gmail.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Filter(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilter_PlusAddressing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plus tag stripped universally",
			input:    "user+tag@hotmail.com",
			expected: "user@hotmail.com",
		},
		{
			name:     "plus tag stripped from outlook",
			input:    "user+newsletter@outlook.com",
			expected: "user@outlook.com",
		},
		{
			name:     "plus tag stripped from protonmail",
			input:    "user+spam@proton.me",
			expected: "user@proton.me",
		},
		{
			name:     "plus with empty tag stripped",
			input:    "user+@gmail.com",
			expected: "user@gmail.com",
		},
		{
			name:     "multiple plus signs keeps first part",
			input:    "user+a+b@gmail.com",
			expected: "user@gmail.com",
		},
		{
			name:     "gmail plus tag and dots both stripped",
			input:    "u.s.e.r+tag@gmail.com",
			expected: "user@gmail.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Filter(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilter_GmailNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "gmail dots removed",
			input:    "user.name@gmail.com",
			expected: "username@gmail.com",
		},
		{
			name:     "gmail multiple dots removed",
			input:    "u.s.e.r@gmail.com",
			expected: "user@gmail.com",
		},
		{
			name:     "googlemail dots removed",
			input:    "user.name@googlemail.com",
			expected: "username@googlemail.com",
		},
		{
			name:     "gmail uppercase local normalized then dots removed",
			input:    "User.Name@Gmail.com",
			expected: "username@gmail.com",
		},
		{
			name:     "non-gmail dots preserved",
			input:    "user.name@outlook.com",
			expected: "user.name@outlook.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Filter(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilter_Blocked(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "whitespace only",
			input: "   ",
		},
		{
			name:  "missing @ separator",
			input: "userexample.com",
		},
		{
			name:  "empty local part",
			input: "@example.com",
		},
		{
			name:  "empty domain",
			input: "user@",
		},
		{
			name:  "local starts with digit",
			input: "1user@example.com",
		},
		{
			name:  "local starts with dot",
			input: ".user@example.com",
		},
		{
			name:  "local starts with plus",
			input: "+user@example.com",
		},
		{
			name:  "consecutive special chars dot-dot",
			input: "user..name@example.com",
		},
		{
			name:  "consecutive special chars dot-plus",
			input: "user.+name@example.com",
		},
		{
			name:  "consecutive special chars minus-underscore",
			input: "user-_name@example.com",
		},
		{
			name:  "invalid character in local part",
			input: "user!name@example.com",
		},
		{
			name:  "local ends with dot",
			input: "user.@example.com",
		},
		{
			name:  "local ends with hyphen",
			input: "user-@example.com",
		},
		{
			name:  "domain with hyphen at start",
			input: "user@-example.com",
		},
		{
			name:  "domain with hyphen at end",
			input: "user@example-.com",
		},
		{
			name:  "domain is bare TLD",
			input: "user@com",
		},
		{
			name:  "invalid character in domain",
			input: "user@exam_ple.com",
		},
		// Government domain blocking
		{
			name:  "government .gov domain",
			input: "user@whitehouse.gov",
		},
		{
			name:  "government .gov.vn domain",
			input: "user@agency.gov.vn",
		},
		{
			name:  "government .gov.uk domain",
			input: "user@service.gov.uk",
		},
		// Plus-only local part becomes empty after stripping
		{
			name:  "plus-only local part",
			input: "+tag@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Filter(tt.input)
			assert.Error(t, err, "expected %q to be blocked", tt.input)
			assert.Empty(t, result)
		})
	}
}

func TestParse(t *testing.T) {
	t.Run("valid email", func(t *testing.T) {
		addr, err := parse("user@example.com")
		require.NoError(t, err)
		assert.Equal(t, "user", addr.Local())
		assert.Equal(t, "example.com", addr.Domain())
		assert.Equal(t, "user@example.com", addr.String())
	})

	t.Run("email with multiple @ signs", func(t *testing.T) {
		addr, err := parse("user@name@example.com")
		require.NoError(t, err)
		// SplitN with 2 splits at first @, so local="user", domain="name@example.com"
		assert.Equal(t, "user", addr.Local())
		assert.Equal(t, "name@example.com", addr.Domain())
	})
}
