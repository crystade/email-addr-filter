# email-addr-filter

This is a Golang library for validating and normalizing email addresses.

## Download

```shell
go get -u github.com/crystade/email-addr-filter
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/crystade/email-addr-filter"
)

func main() {
	// Filter parses, normalizes, and validates the email address
	normalized, err := emailaddr.Filter("  User.Name@Gmail.com  ")
	if err != nil {
		log.Fatalf("Invalid email: %v", err)
	}
	
	fmt.Println(normalized) // Output: username@gmail.com
}
```

## Principles
The library is designed to comply with Crystade best practices on handling email address:
- Normalize-with-evidence: We acknowledge that normalizing email address might run into misidentification, misconfiguration and other unexpected circumstances. Therefore, only supported providers with explicit evidence of email address formats get normalized.
- Block-first: As a universal rule, we block invalid email addresses. They can be lift off on case-by-case basis.

## Rulesets 
Steps: `Universal normalization → Provider-specific normalization → Universal blocker`

### Universal rules
- Normalization:
	+ The local part gets normalized to lowercase
	+ The domain gets normalized to lowercase
	+ Plus addressing is stripped per [RFC 5233](https://datatracker.ietf.org/doc/html/rfc5233) (e.g. `user+tag` → `user`)
- Any violation to a given email address gets blocked:
	+ Characters allowed in the local part are lowercase a-z, 0-9 and special characters `.+-_` (dot, plus, minus and underscore)
	+ The first character in the local part must be lowercase a-z
	+ The last character in the local part must not be a special character
	+ Consecutive special characters in the local part are forbidden, e.g. same `..` or mixed `.+`
	+ The domain must be a valid eTLD+1 (using https://pkg.go.dev/golang.org/x/net/publicsuffix sourced from https://publicsuffix.org/)
	+ Characters allowed in the hostname (after lowercasing) must be a-z, 0-9 and `-` (hyphen)
	+ The first character in the hostname must not be a hyphen `-`
	+ The last character in the hostname must not be a hyphen `-`
	+ IANA reserved domains (e.g. `example.com`, `example.org`) are blocked
	+ Disposable/temporary email domains are blocked (configurable via `blocked_disposable_domains.txt`)
	+ Government email domains are blocked: `.gov`, `*.gov.*` patterns are blocked by code, with additional suffixes configurable via `blocked_government_domains.txt`

```
addr = local-part "@" domain
    (username)

domain = hostname { "." hostname }
```

### Provider-specific rules

#### Gmail
- Normalization: All dots `.` are removed (https://support.google.com/mail/answer/7436150)

### Configurable blocklists

The library embeds two text files for domain blocking. Each file contains one domain (or suffix) per line, with `#` comments and empty lines ignored.

| File | Purpose |
|---|---|
| `blocked_disposable_domains.txt` | Exact-match disposable/temporary email domains |
| `blocked_government_domains.txt` | Additional government domain suffixes (`.gov`/`*.gov.*` are already handled in code) |

## References
- [RFC 5233: Sieve Email Filtering: Subaddress Extension](https://datatracker.ietf.org/doc/html/rfc5233)
- [RFC 5322: Internet Message Format](https://datatracker.ietf.org/doc/html/rfc5322)
- [Salesforce Email Address Format Technical Standards Validation](https://help.salesforce.com/s/articleView?id=000384328)
- https://stackoverflow.com/a/2049510
- https://github.com/disposable-email-domains/disposable-email-domains

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
