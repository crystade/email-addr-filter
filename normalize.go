package emailaddr

import "strings"

// normalizeUniversal applies universal normalization rules:
//   - Local part gets lowercased
//   - Domain gets lowercased
//   - Plus addressing is stripped (RFC 5233)
func (a *Addr) normalizeUniversal() {
	a.local = strings.ToLower(a.local)
	a.domain = strings.ToLower(a.domain)

	// Strip plus addressing (RFC 5233: Sieve Email Filtering: Subaddress Extension)
	// e.g. "user+tag" → "user"
	if idx := strings.IndexByte(a.local, '+'); idx >= 0 {
		a.local = a.local[:idx]
	}
}

// normalizeProvider applies provider-specific normalization rules
// based on the domain of the email address.
func (a *Addr) normalizeProvider() {
	host := hostname(a.domain)
	switch host {
	case "gmail.com", "googlemail.com":
		// Gmail: All dots in the local part are removed
		// https://support.google.com/mail/answer/7436150
		a.local = strings.ReplaceAll(a.local, ".", "")
	}
}

// hostname returns the full lowercased domain for provider matching.
// Since normalizeUniversal already lowercases the domain before this
// is called, this function simply returns it as-is.
func hostname(domain string) string {
	return domain
}
