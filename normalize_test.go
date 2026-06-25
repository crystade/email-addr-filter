package emailaddr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeUniversal(t *testing.T) {
	tests := []struct {
		name           string
		local          string
		domain         string
		expectedLocal  string
		expectedDomain string
	}{
		{
			name:           "lowercases domain",
			local:          "user",
			domain:         "EXAMPLE.COM",
			expectedLocal:  "user",
			expectedDomain: "example.com",
		},
		{
			name:           "lowercases local part",
			local:          "USER",
			domain:         "example.com",
			expectedLocal:  "user",
			expectedDomain: "example.com",
		},
		{
			name:           "lowercases both",
			local:          "UsEr",
			domain:         "ExAmPlE.CoM",
			expectedLocal:  "user",
			expectedDomain: "example.com",
		},
		{
			name:           "already lowercase is unchanged",
			local:          "user",
			domain:         "example.com",
			expectedLocal:  "user",
			expectedDomain: "example.com",
		},
		{
			name:           "strips plus addressing",
			local:          "user+tag",
			domain:         "example.com",
			expectedLocal:  "user",
			expectedDomain: "example.com",
		},
		{
			name:           "strips plus with empty tag",
			local:          "user+",
			domain:         "example.com",
			expectedLocal:  "user",
			expectedDomain: "example.com",
		},
		{
			name:           "no plus sign unchanged",
			local:          "username",
			domain:         "example.com",
			expectedLocal:  "username",
			expectedDomain: "example.com",
		},
		{
			name:           "plus stripping with uppercase",
			local:          "User+Tag",
			domain:         "example.com",
			expectedLocal:  "user",
			expectedDomain: "example.com",
		},
		{
			name:           "multiple plus signs keeps first part only",
			local:          "user+tag1+tag2",
			domain:         "example.com",
			expectedLocal:  "user",
			expectedDomain: "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := &Addr{local: tt.local, domain: tt.domain}
			addr.normalizeUniversal()
			assert.Equal(t, tt.expectedLocal, addr.local)
			assert.Equal(t, tt.expectedDomain, addr.domain)
		})
	}
}

func TestNormalizeProvider(t *testing.T) {
	tests := []struct {
		name          string
		local         string
		domain        string
		expectedLocal string
	}{
		{
			name:          "gmail removes dots",
			local:         "user.name",
			domain:        "gmail.com",
			expectedLocal: "username",
		},
		{
			name:          "gmail removes multiple dots",
			local:         "u.s.e.r",
			domain:        "gmail.com",
			expectedLocal: "user",
		},
		{
			name:          "googlemail removes dots",
			local:         "user.name",
			domain:        "googlemail.com",
			expectedLocal: "username",
		},
		{
			name:          "non-gmail preserves dots",
			local:         "user.name",
			domain:        "outlook.com",
			expectedLocal: "user.name",
		},
		{
			name:          "gmail no dots unchanged",
			local:         "username",
			domain:        "gmail.com",
			expectedLocal: "username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := &Addr{local: tt.local, domain: tt.domain}
			addr.normalizeProvider()
			assert.Equal(t, tt.expectedLocal, addr.local)
		})
	}
}

func TestHostname(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected string
	}{
		{"simple domain", "gmail.com", "gmail.com"},
		{"subdomain", "mail.google.com", "mail.google.com"},
		{"already lowercase", "example.com", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, hostname(tt.domain))
		})
	}
}
