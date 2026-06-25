package emailaddr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDisposableDomain(t *testing.T) {
	// With an empty blocklist, nothing should be blocked
	t.Run("empty list allows all domains", func(t *testing.T) {
		assert.True(t, isDisposableDomain("mailinator.com"))
	})
}

func TestIsGovernmentDomain(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		want   bool
	}{
		// Code-level .gov pattern
		{"bare .gov TLD", "whitehouse.gov", true},
		{"subdomain of .gov", "mail.whitehouse.gov", true},

		// Code-level *.gov.* pattern
		{".gov.vn", "agency.gov.vn", true},
		{".gov.uk", "service.gov.uk", true},
		{".gov.au", "tax.gov.au", true},
		{"subdomain of .gov.vn", "sub.agency.gov.vn", true},

		// Non-government domains
		{"normal .com", "example.com", false},
		{"govtrack.us not blocked", "govtrack.us", false},
		{"my-gov.com not blocked", "my-gov.com", false},
		{"gmail.com not blocked", "gmail.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isGovernmentDomain(tt.domain))
		})
	}
}

func TestParseLines(t *testing.T) {
	t.Run("ignores comments and blanks", func(t *testing.T) {
		data := "# comment\n\nmailinator.com\n  yopmail.com  \n# another comment\n"
		result := parseLines(data)
		assert.Len(t, result, 2)
		_, ok1 := result["mailinator.com"]
		_, ok2 := result["yopmail.com"]
		assert.True(t, ok1)
		assert.True(t, ok2)
	})

	t.Run("lowercases entries", func(t *testing.T) {
		data := "Mailinator.COM\n"
		result := parseLines(data)
		_, ok := result["mailinator.com"]
		assert.True(t, ok)
	})

	t.Run("empty data", func(t *testing.T) {
		result := parseLines("")
		assert.Empty(t, result)
	})
}

func TestParseSuffixes(t *testing.T) {
	t.Run("adds leading dot if missing", func(t *testing.T) {
		data := "mil\n.go.jp\n"
		result := parseSuffixes(data)
		assert.Equal(t, []string{".mil", ".go.jp"}, result)
	})

	t.Run("ignores comments and blanks", func(t *testing.T) {
		data := "# comment\n\n.mil\n"
		result := parseSuffixes(data)
		assert.Equal(t, []string{".mil"}, result)
	})

	t.Run("empty data", func(t *testing.T) {
		result := parseSuffixes("")
		assert.Nil(t, result)
	})
}
