package emailaddr

import (
	"fmt"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// validate checks the universal blocker rules.
// Any violation results in the address being blocked.
func (a *Addr) validate() error {
	if err := validateLocal(a.local); err != nil {
		return err
	}
	if err := validateDomain(a.domain); err != nil {
		return err
	}
	return nil
}

// validateLocal checks the local part against the universal blocker rules:
//   - Characters allowed: lowercase a-z, 0-9, and . + - _
//   - First character must be lowercase a-z
//   - Consecutive special characters are forbidden
func validateLocal(local string) error {
	if local == "" {
		return fmt.Errorf("empty local part")
	}

	// First character must be lowercase a-z
	if local[0] < 'a' || local[0] > 'z' {
		return fmt.Errorf("local part must start with a lowercase letter")
	}

	specials := ".+-_"

	for i := 0; i < len(local); i++ {
		ch := local[i]

		// Allowed: a-z, 0-9, and . + - _
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			continue
		}
		if strings.ContainsRune(specials, rune(ch)) {
			// Check for consecutive special characters
			if i > 0 && strings.ContainsRune(specials, rune(local[i-1])) {
				return fmt.Errorf("consecutive special characters not allowed in local part")
			}
			continue
		}
		return fmt.Errorf("invalid character %q in local part", ch)
	}

	// Last character must not be a special character
	if strings.ContainsRune(specials, rune(local[len(local)-1])) {
		return fmt.Errorf("local part must not end with a special character")
	}

	return nil
}

// validateDomain checks the domain against the universal blocker rules:
//   - Domain must be eTLD+1 (valid public suffix)
//   - Characters allowed in hostname: a-z, 0-9, -
//   - First character of hostname must not be a hyphen
func validateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("empty domain")
	}

	// Must be a valid eTLD+1
	if err := validateETLDPlusOne(domain); err != nil {
		return err
	}

	// Top-level domain must be a recognized ICANN TLD
	if err := validateTLD(domain); err != nil {
		return err
	}

	// Validate each hostname component
	labels := strings.Split(domain, ".")
	for _, label := range labels {
		if err := validateHostname(label); err != nil {
			return fmt.Errorf("invalid hostname %q: %w", label, err)
		}
	}

	// Block IANA reserved domains
	if isIANAReservedDomain(domain) {
		return fmt.Errorf("IANA reserved domain %q is not allowed", domain)
	}

	// Block disposable/temporary email domains
	if isDisposableDomain(domain) {
		return fmt.Errorf("disposable email domain %q is not allowed", domain)
	}

	// Block government domains to avoid legal issues
	if isGovernmentDomain(domain) {
		return fmt.Errorf("government email domain %q is not allowed", domain)
	}

	return nil
}

// validateETLDPlusOne checks that the domain is at least eTLD+1 using
// the Public Suffix List.
func validateETLDPlusOne(domain string) error {
	// publicsuffix.EffectiveTLDPlusOne returns the eTLD+1 for a domain.
	// If the domain itself is just a public suffix (eTLD), it returns
	// an empty string or errors.
	etld1, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return fmt.Errorf("domain %q is not a valid eTLD+1: %w", domain, err)
	}
	if etld1 == "" {
		return fmt.Errorf("domain %q is not a valid eTLD+1: bare public suffix", domain)
	}
	return nil
}

// validateTLD checks that the rightmost label (top-level domain) of a domain
// is a recognized ICANN-managed TLD. This prevents domains with completely
// fabricated TLDs like "com8" from passing validation, as the publicsuffix
// package's wildcard rule would otherwise treat any unknown TLD as valid.
func validateTLD(domain string) error {
	labels := strings.Split(domain, ".")
	tld := labels[len(labels)-1]
	_, icann := publicsuffix.PublicSuffix(tld)
	if !icann {
		return fmt.Errorf("top-level domain %q is not a recognized ICANN TLD", tld)
	}
	return nil
}

// validateHostname checks a single hostname label against the rules:
//   - Characters allowed: a-z, 0-9, -
//   - First character must not be a hyphen
func validateHostname(label string) error {
	if label == "" {
		return fmt.Errorf("empty hostname label")
	}

	// First character must not be a hyphen
	if label[0] == '-' {
		return fmt.Errorf("hostname must not start with a hyphen")
	}

	// Last character must not be a hyphen
	if label[len(label)-1] == '-' {
		return fmt.Errorf("hostname must not end with a hyphen")
	}

	for i := 0; i < len(label); i++ {
		ch := label[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			continue
		}
		return fmt.Errorf("invalid character %q in hostname", ch)
	}

	return nil
}
