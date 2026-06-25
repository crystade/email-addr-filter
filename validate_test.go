package emailaddr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLocal(t *testing.T) {
	tests := []struct {
		name    string
		local   string
		wantErr bool
	}{
		{"valid simple", "user", false},
		{"valid with dot", "user.name", false},
		{"valid with plus", "user+tag", false},
		{"valid with hyphen", "user-name", false},
		{"valid with underscore", "user_name", false},
		{"valid with digits", "user123", false},
		{"valid complex", "a.b+c-d_e", false},

		{"empty", "", true},
		{"starts with digit", "1user", true},
		{"starts with dot", ".user", true},
		{"starts with hyphen", "-user", true},
		{"starts with plus", "+user", true},
		{"starts with underscore", "_user", true},
		{"consecutive dots", "user..name", true},
		{"consecutive mixed specials", "user.-name", true},
		{"invalid char space", "user name", true},
		{"invalid char exclamation", "user!name", true},
		{"invalid char at", "user@name", true},
		{"invalid char hash", "user#name", true},
		{"ends with dot", "user.", true},
		{"ends with hyphen", "user-", true},
		{"ends with plus", "user+", true},
		{"ends with underscore", "user_", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLocal(tt.local)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr bool
	}{
		{"valid simple", "my-domain.org", false},
		{"valid subdomain", "mail.my-domain.org", false},
		{"valid with hyphen", "my-domain.com", false},
		{"valid with digits", "example123.com", false},

		{"empty", "", true},
		{"bare TLD", "com", true},
		{"starts with hyphen", "-my-domain.com", true},
		{"label starts with hyphen", "my-domain.-com.org", true},
		{"invalid char underscore", "my_domain.com", true},
		{"empty label (double dot)", "my-domain..com", true},
		{"label ends with hyphen", "my-domain-.com", true},

		// IANA reserved domains
		{"IANA example.com", "example.com", true},
		{"IANA example.org subdomain", "mail.example.org", true},
		{"IANA example.net", "example.net", true},
		{"IANA example.edu", "example.edu", true},

		// Disposable domain blocking
		{"disposable 0-mail.com", "0-mail.com", true},
		{"disposable mailinator.com", "mailinator.com", true},

		// Government domain blocking
		{"government .gov", "whitehouse.gov", true},
		{"government .gov.vn", "agency.gov.vn", true},
		{"government .gov.uk", "service.gov.uk", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDomain(tt.domain)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateHostname(t *testing.T) {
	tests := []struct {
		name    string
		label   string
		wantErr bool
	}{
		{"valid lowercase", "example", false},
		{"valid with digits", "example123", false},
		{"valid with hyphen", "my-example", false},

		{"empty", "", true},
		{"starts with hyphen", "-example", true},
		{"ends with hyphen", "example-", true},
		{"uppercase letter", "Example", true},
		{"underscore", "my_example", true},
		{"space", "my example", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHostname(tt.label)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
