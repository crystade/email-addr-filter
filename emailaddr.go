// Package emailaddr provides email address validation and normalization
// according to Crystade best practices.
//
// Principles:
//   - Normalize-with-evidence: Only supported providers with explicit evidence
//     of email address formats get normalized.
//   - Block-first: Invalid email addresses are blocked by default.
//
// Processing pipeline: Universal normalization → Provider-specific normalization → Universal blocker
package emailaddr

import (
	"fmt"
	"strings"
)

// Addr represents a parsed email address with separate local part and domain.
type Addr struct {
	local  string
	domain string
}

// Local returns the local part of the email address (before the @).
func (a *Addr) Local() string { return a.local }

// Domain returns the domain part of the email address (after the @).
func (a *Addr) Domain() string { return a.domain }

// String returns the email address as "local@domain".
func (a *Addr) String() string {
	return a.local + "@" + a.domain
}

// Filter parses, normalizes, and validates an email address.
// It returns the normalized email address string or an error if the
// address is invalid and should be blocked.
func Filter(raw string) (string, error) {
	addr, err := parse(raw)
	if err != nil {
		return "", err
	}

	// Step 1: Universal normalization
	addr.normalizeUniversal()

	// Step 2: Provider-specific normalization
	addr.normalizeProvider()

	// Step 3: Universal blocker validation
	if err := addr.validate(); err != nil {
		return "", fmt.Errorf("blocked: %w", err)
	}

	return addr.String(), nil
}

// parse splits a raw email address into local part and domain.
func parse(raw string) (*Addr, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty email address")
	}

	parts := strings.SplitN(raw, "@", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("missing @ separator")
	}

	local := parts[0]
	domain := parts[1]

	if local == "" {
		return nil, fmt.Errorf("empty local part")
	}
	if domain == "" {
		return nil, fmt.Errorf("empty domain")
	}

	return &Addr{local: local, domain: domain}, nil
}
