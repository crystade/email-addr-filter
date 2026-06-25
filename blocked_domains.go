package emailaddr

import (
	_ "embed"
	"strings"
)

//go:embed blocked_disposable_domains.txt
var disposableDomainsData string

//go:embed blocked_government_domains.txt
var governmentDomainsData string

var (
	disposableDomains  map[string]struct{}
	governmentSuffixes []string
)

func init() {
	disposableDomains = parseLines(disposableDomainsData)
	governmentSuffixes = parseSuffixes(governmentDomainsData)
}

// parseLines parses a newline-separated list into a set,
// ignoring empty lines and lines starting with #.
func parseLines(data string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		m[strings.ToLower(line)] = struct{}{}
	}
	return m
}

// parseSuffixes parses a newline-separated list of domain suffixes,
// ensuring each entry starts with a dot.
func parseSuffixes(data string) []string {
	var suffixes []string
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		s := strings.ToLower(line)
		if !strings.HasPrefix(s, ".") {
			s = "." + s
		}
		suffixes = append(suffixes, s)
	}
	return suffixes
}

// isDisposableDomain checks if the domain is a known disposable email domain
// by exact match against the embedded blocklist.
func isDisposableDomain(domain string) bool {
	_, ok := disposableDomains[domain]
	return ok
}

// isGovernmentDomain checks if the domain is a government domain.
// It matches:
//   - Domains ending with ".gov" (e.g. whitehouse.gov)
//   - Domains containing ".gov." as a segment (e.g. agency.gov.vn)
//   - Domains matching additional suffixes from blocked_government_domains.txt
func isGovernmentDomain(domain string) bool {
	// Code-level pattern: .gov and *.gov.*
	if strings.HasSuffix(domain, ".gov") || strings.Contains(domain, ".gov.") {
		return true
	}

	// Additional suffixes from the embedded list
	for _, suffix := range governmentSuffixes {
		if strings.HasSuffix(domain, suffix) {
			return true
		}
	}
	return false
}

// isIANAReservedDomain checks if the domain is an IANA reserved domain
// or a subdomain thereof (e.g. example.com, example.org, example.net, example.edu).
func isIANAReservedDomain(domain string) bool {
	if domain == "example.com" || strings.HasSuffix(domain, ".example.com") ||
		domain == "example.org" || strings.HasSuffix(domain, ".example.org") ||
		domain == "example.net" || strings.HasSuffix(domain, ".example.net") ||
		domain == "example.edu" || strings.HasSuffix(domain, ".example.edu") {
		return true
	}
	return false
}
