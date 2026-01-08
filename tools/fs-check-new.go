package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/miekg/dns"
)

// fs-check: sanity checker for _for-sale DNS TXT records
//
// Usage: fs-check domain.tld
//
// Flags:
//   -json                output machine-readable JSON; JSON output includes full values
//
// Behavior:
//   - queries resolver(s) from /etc/resolv.conf using EDNS0 with larger UDP buffer
//     and falls back to TCP if the UDP reply is truncated.
//   - validates TXT RRs at _for-sale.<domain> according to draft-davids-forsalereg-19,
//     including UTF-8 / control-character checks derived from the draft's
//     recommendation about encoding and Unicode subsets.
//   - decodes presentation escapes (e.g., \240\159\142\133) into raw bytes before parsing
//   - prints human-readable diagnostics or JSON (when -json is set)
//   - output is sorted: VALID, INVALID, IGNORED (both human and JSON modes)
//
// Exit codes:
//   0 : at least one valid _for-sale TXT record found (with or without warnings)
//   2 : TXT records found but none considered valid (all invalid/ignored)
//   3 : usage error or DNS/network error

const (
	versionTag     = "v=FORSALE1;"
	maxTxtStr      = 255
	defaultTimeout = 5 * time.Second
)

var (
	// fval: currency (one or more uppercase letters) followed by amount (digits, optional .fraction)
	fvalRe = regexp.MustCompile(`^[A-Z]+[0-9]+(?:\.[0-9]+)?$`)
)

// recordResult holds diagnostics for one TXT RR. `Rr` is omitted from JSON.
type recordResult struct {
	Rr                     dns.RR   `json:"-"`
	Content                string   `json:"content"`                       // decoded concatenated character-strings (full)
	RawTxts                []string `json:"raw_txts"`                      // raw presentation strings from the RR
	RawDecodedLens         []int    `json:"raw_decoded_lens,omitempty"`    // decoded octet length per raw part
	TTL                    uint32   `json:"ttl"`
	RawCount               int      `json:"raw_count"`                     // number of character-strings as seen in the RR
	FitsSingleCharstring   bool     `json:"fits_single_charstring"`        // true if concatenation fits in single char-string <=255
	Valid                  bool     `json:"valid"`
	Ignored                bool     `json:"ignored"`
	Tag                    string   `json:"tag,omitempty"`
	TagValue               string   `json:"tag_value,omitempty"`
	Messages               []string `json:"messages,omitempty"`
	ConcatenatedLength     int      `json:"concatenated_length"` // bytes
}

type jsonOutput struct {
	Query        string         `json:"query"`
	Records      []recordResult `json:"records"`
	TTLCounts    map[uint32]int `json:"ttl_counts,omitempty"`
	Summary      string         `json:"summary"`
	ValidCount   int            `json:"valid_count"`
	IgnoredCount int            `json:"ignored_count"`
	InvalidCount int            `json:"invalid_count"`
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] domain\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Example: %s example.com\n", os.Args[0])
		flag.PrintDefaults()
	}
	jsonOutFlag := flag.Bool("json", false, "output machine-readable JSON (includes full values)")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Error: missing domain argument.")
		flag.Usage()
		os.Exit(3)
	}
	domain := strings.TrimSpace(flag.Arg(0))
	if domain == "" {
		fmt.Fprintln(os.Stderr, "Error: empty domain.")
		os.Exit(3)
	}

	lower := strings.ToLower(strings.TrimSuffix(domain, "."))
	if lower == "arpa" || strings.HasSuffix(lower, ".arpa") {
		fmt.Printf("Domain %q is in the .arpa hierarchy - records under .arpa are out of scope and MUST be ignored per the draft.\n", domain)
		os.Exit(0)
	}

	fqdn := dns.Fqdn("_for-sale." + domain) // trailing dot

	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read /etc/resolv.conf: %v\n", err)
		os.Exit(3)
	}
	if len(conf.Servers) == 0 {
		fmt.Fprintf(os.Stderr, "No DNS servers configured in resolv.conf\n")
		os.Exit(3)
	}

	// Query using EDNS0 and TCP fallback when truncated (fixed default timeout)
	var resp *dns.Msg
	var lastErr error
	for _, server := range conf.Servers {
		serverAddr := server + ":" + conf.Port
		r, err := queryTXTWithFallback(fqdn, serverAddr, defaultTimeout)
		if err != nil {
			lastErr = err
			continue
		}
		resp = r
		break
	}
	if resp == nil {
		fmt.Fprintf(os.Stderr, "DNS query failed: %v\n", lastErr)
		os.Exit(3)
	}

	// collect TXT answers
	var txtRRs []*dns.TXT
	for _, a := range resp.Answer {
		if t, ok := a.(*dns.TXT); ok {
			txtRRs = append(txtRRs, t)
		}
	}

	if len(txtRRs) == 0 {
		fmt.Printf("No TXT records found at %s\n", fqdn)
		os.Exit(2)
	}

	results := make([]recordResult, 0, len(txtRRs))
	ttls := make(map[uint32]int)
	var anyValid bool

	for _, t := range txtRRs {
		res := analyzeTXT(t)
		res.Rr = t
		results = append(results, res)
		if res.TTL > 0 {
			ttls[res.TTL]++
		}
		if res.Valid {
			anyValid = true
		}
	}

	// Sort results into groups: VALID, INVALID, IGNORED (preserve order within group)
	valids := make([]recordResult, 0)
	invalids := make([]recordResult, 0)
	ignored := make([]recordResult, 0)
	for _, r := range results {
		if r.Ignored {
			ignored = append(ignored, r)
		} else if r.Valid {
			valids = append(valids, r)
		} else {
			invalids = append(invalids, r)
		}
	}
	sorted := make([]recordResult, 0, len(results))
	sorted = append(sorted, valids...)
	sorted = append(sorted, invalids...)
	sorted = append(sorted, ignored...)

	// JSON mode: emit structured output including full values (no truncation)
	if *jsonOutFlag {
		out := jsonOutput{
			Query:     fqdn,
			Records:   sorted,
			TTLCounts: ttls,
		}
		validCount := len(valids)
		ignoredCount := len(ignored)
		invalidCount := len(invalids)
		out.ValidCount = validCount
		out.IgnoredCount = ignoredCount
		out.InvalidCount = invalidCount
		out.Summary = fmt.Sprintf("%d record(s) total: %d valid, %d ignored (no version), %d invalid",
			len(sorted), validCount, ignoredCount, invalidCount)

		enc, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to marshal JSON output: %v\n", err)
			os.Exit(3)
		}
		fmt.Println(string(enc))
		// exit code based on validity
		if anyValid {
			os.Exit(0)
		}
		os.Exit(2)
	}

	// Human-readable output (always full content), sorted as requested
	fmt.Printf("Found %d TXT record(s) at %s\n\n", len(sorted), fqdn)
	for i, r := range sorted {
		fmt.Printf("Record #%d (TTL=%d, raw-strings=%d, concatenated-bytes=%d, fits_single_charstring=%v):\n",
			i+1, r.TTL, r.RawCount, r.ConcatenatedLength, r.FitsSingleCharstring)
		// show raw character-strings as returned by miekg/dns (presentation form) for diagnostic.
		// To mimic dig output we avoid fmt %q (which produces Go-style escaping). Print the presentation string inside quotes,
		// but show the decoded octet length (not the presentation length).
		if r.Rr != nil {
			if t, ok := r.Rr.(*dns.TXT); ok {
				for si, s := range t.Txt {
					decodedLen := 0
					if si < len(r.RawDecodedLens) {
						decodedLen = r.RawDecodedLens[si]
					}
					fmt.Printf("  Raw Txt[%d] (decoded-len=%d): \"%s\"\n", si, decodedLen, s)
				}
			}
		}
		// show decoded content (always full)
		if r.Content != "" {
			fmt.Printf("  Decoded content (len=%d): %s\n", len(r.Content), r.Content)
		} else {
			fmt.Printf("  Decoded content: <empty>\n")
		}
		if r.Ignored {
			fmt.Printf("  Verdict: IGNORED (no valid version tag found)\n")
		} else if r.Valid {
			fmt.Printf("  Verdict: VALID\n")
		} else {
			fmt.Printf("  Verdict: INVALID\n")
		}
		if r.Tag != "" {
			fmt.Printf("  Content tag: %s\n", r.Tag)
			if r.TagValue != "" {
				fmt.Printf("  Content value (len=%d): %s\n", len(r.TagValue), r.TagValue)
			}
		} else {
			fmt.Printf("  No content tag present (empty content after version tag)\n")
		}
		for _, m := range r.Messages {
			fmt.Printf("  - %s\n", m)
		}
		fmt.Println()
	}

	// RRset TTL checks
	if len(ttls) > 1 {
		fmt.Printf("Warning: TXT RRset contains records with differing TTLs (RRset TTLs must be the same per RFC2181 Section 5.2). TTLs seen:\n")
		// Print TTLs in deterministic order
		keys := make([]int, 0, len(ttls))
		for k := range ttls {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)
		for _, ttl := range keys {
			fmt.Printf("  TTL %d: %d record(s)\n", ttl, ttls[uint32(ttl)])
		}
	} else {
		for ttl := range ttls {
			if ttl > 3600 {
				fmt.Printf("Warning: TTL=%d is greater than the recommended 3600s (1 hour). Long TTLs increase the risk of outdated sale information.\n", ttl)
			}
		}
	}

	// Summary & exit code
	validCount := len(valids)
	ignoredCount := len(ignored)
	invalidCount := len(invalids)

	fmt.Printf("\nSummary: %d record(s) total: %d valid, %d ignored (no version), %d invalid\n",
		len(sorted), validCount, ignoredCount, invalidCount)

	if anyValid {
		os.Exit(0)
	}
	os.Exit(2)
}

// queryTXTWithFallback performs a TXT query for qname to serverAddr.
// It sets EDNS0 with a larger UDP buffer (4096) and, if the response is truncated,
// retries the same query over TCP to obtain the full response.
func queryTXTWithFallback(qname, serverAddr string, timeout time.Duration) (*dns.Msg, error) {
	client := &dns.Client{Timeout: timeout, Net: "udp"}
	msg := new(dns.Msg)
	msg.SetQuestion(qname, dns.TypeTXT)
	// request EDNS0 with larger UDP payload (4096)
	msg.SetEdns0(4096, false)

	r, _, err := client.Exchange(msg, serverAddr)
	if err != nil {
		return nil, err
	}
	if r.Truncated {
		// Retry over TCP to obtain full answer
		client.Net = "tcp"
		r2, _, err2 := client.Exchange(msg, serverAddr)
		if err2 != nil {
			// return the original truncated response if TCP failed, but signal error
			return r, fmt.Errorf("UDP response truncated; TCP retry failed: %w", err2)
		}
		return r2, nil
	}
	return r, nil
}

// analyzeTXT validates a dns.TXT RR and returns a recordResult.
// Improvements over previous logic:
//  - Always concatenate unescaped character-strings to form the logical content.
//  - Provide a heuristic `FitsSingleCharstring` that is true when concatenated content <=255
//    allowing the tool to treat some multi-part TXT RRs as a single logical string for robustness.
//  - Keep and report the raw count but do not automatically reject records simply because
//    the server split a logical single string into multiple character-strings (common in practice).
//  - Validate textual content for UTF-8 and forbidden control characters per the draft's recommendations.
//  - Add a per-record TTL warning when the observed TTL exceeds the recommended 3600 seconds
//    (noting that this may come from a resolver cache).
//  - Add an explicit warning when TXT RRs are multi-part (raw_count>1) and report decoded per-part lengths.
func analyzeTXT(t *dns.TXT) recordResult {
	res := recordResult{
		TTL:     t.Hdr.Ttl,
		RawTxts: append([]string(nil), t.Txt...),
		RawCount: len(t.Txt),
	}

	// Warn when observed TTL exceeds recommended value (note: may be a cached reply)
	if res.TTL > 3600 {
		res.Messages = append(res.Messages, fmt.Sprintf("Warning: observed TTL=%d exceeds the recommended 3600s (1 hour). Note: this value may come from a resolver cache and not the authoritative server.", res.TTL))
	}

	// If zero character-strings, invalid
	if res.RawCount == 0 {
		res.Messages = append(res.Messages, "TXT RR contains zero character-strings (invalid).")
		res.Valid = false
		return res
	}

	// Decode each character-string from presentation escapes to raw bytes, collect decoded lengths, then concatenate
	var b []byte
	decodedLens := make([]int, 0, len(t.Txt))
	for _, part := range t.Txt {
		ub, err := unescapePresentation(part)
		if err != nil {
			// record the error but continue; ub may contain partial decoded bytes
			res.Messages = append(res.Messages, fmt.Sprintf("Warning: error while unescaping presentation string: %v", err))
		}
		decodedLens = append(decodedLens, len(ub))
		b = append(b, ub...)
	}
	res.RawDecodedLens = decodedLens

	res.ConcatenatedLength = len(b)
	res.Content = string(b)

	// Heuristic: consider the concatenation as a single logical string if concatenated length <= 255
	// This handles cases where servers split long quoted string data into multiple character-strings.
	if res.ConcatenatedLength <= maxTxtStr {
		res.FitsSingleCharstring = true
	} else {
		res.FitsSingleCharstring = false
	}

	// Warn if the server split into multiple raw character-strings
	if res.RawCount > 1 {
		// explicit multi-part warning requested by user
		res.Messages = append(res.Messages, fmt.Sprintf("Warning: TXT RR contains %d raw character-strings (multi-part RR). The draft RECOMMENDS using a single character-string; consider converting to a single string to avoid ambiguity.", res.RawCount))
		// If any raw part decoded length >255 (octets), it's definitely non-conformant
		for i, decLen := range decodedLens {
			if decLen > maxTxtStr {
				res.Messages = append(res.Messages, fmt.Sprintf("Raw character-string #%d decoded length %d octets exceeds 255 octets (maximum).", i, decLen))
				// mark invalid, but continue diagnostics
				res.Valid = false
			}
		}
	}

	// If concatenation exceeds 255 bytes, that's non-conformant with the draft's requirement that
	// each TXT record's RDATA MUST be a single character-string of at most 255 bytes.
	if res.ConcatenatedLength > maxTxtStr {
		res.Messages = append(res.Messages, fmt.Sprintf("Decoded concatenated content byte length %d exceeds 255 octets; this is non-conformant.", res.ConcatenatedLength))
		// keep going to give diagnostics, but mark invalid
		res.Valid = false
	}

	// Check version tag presence (must be at start, case-sensitive) on the concatenated content
	content := res.Content
	if !strings.HasPrefix(content, versionTag) {
		// Robustness: accept versionTag followed by single space or tab (warn)
		if strings.HasPrefix(content, versionTag+" ") || strings.HasPrefix(content, versionTag+"\t") {
			res.Messages = append(res.Messages, "Record starts with version tag followed by whitespace - accepted under robustness, but spaces are not allowed by the ABNF (warning).")
			rest := content[len(versionTag):]
			if len(rest) > 0 && (rest[0] == ' ' || rest[0] == '\t') {
				rest = rest[1:]
			}
			content = rest
		} else {
			res.Messages = append(res.Messages, "No valid version tag found at start of the TXT record. TXT records without the exact, case-sensitive version tag \"v=FORSALE1;\" MUST NOT be interpreted as valid _for-sale indicators (this record will be ignored).")
			res.Ignored = true
			res.Valid = false
			return res
		}
	} else {
		content = content[len(versionTag):]
	}

	// If no content after version tag => valid indicator with no further info
	if content == "" {
		res.Valid = true
		res.Messages = append(res.Messages, "Record contains only the version tag and no content: valid indicator that the domain is for sale.")
		return res
	}

	// Content must be exactly one tag-value pair
	tags := []string{"fcod=", "ftxt=", "furi=", "fval="}
	foundTag := ""
	for _, tg := range tags {
		if strings.HasPrefix(content, tg) {
			foundTag = tg[:len(tg)-1] // remove '='
			break
		}
	}
	if foundTag == "" {
		res.Messages = append(res.Messages, fmt.Sprintf("Content does not start with a recognised content tag (fcod=, ftxt=, furi=, fval=). Found content: %q", content))
		// Per draft: if version present but content invalid, processors SHOULD assume domain is for sale
		res.Valid = true
		res.Messages = append(res.Messages, "Per the draft, since a valid version tag is present but content is invalid, processors SHOULD still treat the domain as for sale. This tool marks the record as ACCEPTED (but with warnings).")
		return res
	}

	val := content[len(foundTag)+1:]
	res.Tag = foundTag
	res.TagValue = val

	// detect ambiguous constructs: additional tag markers inside value
	for _, tg := range tags {
		needle := ";" + tg
		if strings.Contains(val, needle) {
			res.Messages = append(res.Messages, fmt.Sprintf("The content value contains %q which looks like an additional tag-value pair. The draft REQUIRES exactly one tag-value pair per record; embedding additional tags in the value can be ambiguous.", needle))
		}
	}

	// Validate per tag
	switch foundTag {
	case "fcod":
		if len([]byte(val)) < 1 {
			res.Messages = append(res.Messages, "fcod= has an empty value (must be at least 1 octet).")
			res.Valid = false
			return res
		}
		if len([]byte(val)) > 239 {
			res.Messages = append(res.Messages, fmt.Sprintf("fcod value byte length %d exceeds the draft's maximum of 239 octets for fcod-value.", len([]byte(val))))
			res.Valid = false
			return res
		}
		// fcod is opaque; do not apply UTF-8 checks
		res.Valid = true
		res.Messages = append(res.Messages, "fcod content tag is syntactically acceptable (semantic interpretation is proprietary).")
		return res

	case "ftxt":
		// ftxt-value = 1*239OCTET
		if len([]byte(val)) < 1 {
			res.Messages = append(res.Messages, "ftxt= has an empty value (must be at least 1 octet).")
			res.Valid = false
			return res
		}
		if len([]byte(val)) > 239 {
			res.Messages = append(res.Messages, fmt.Sprintf("ftxt value byte length %d exceeds the draft's maximum of 239 octets for ftxt-value.", len([]byte(val))))
			res.Valid = false
			return res
		}

		// New: enforce recommendations about UTF-8 / control characters:
		warns, errs := checkUnicodeContent(val)
		for _, w := range warns {
			res.Messages = append(res.Messages, "Warning: "+w)
		}
		if len(errs) > 0 {
			for _, e := range errs {
				res.Messages = append(res.Messages, "Error: "+e)
			}
			res.Valid = false
			return res
		}

		res.Valid = true
		res.Messages = append(res.Messages, "ftxt content tag is syntactically acceptable. Note: avoid using URIs in ftxt; prefer furi=. Ensure non-ASCII text is UTF-8 encoded.")
		return res

	case "furi":
		if len(val) < 1 {
			res.Messages = append(res.Messages, "furi= has an empty value (must contain exactly one URI or IRI).")
			res.Valid = false
			return res
		}

		// Check for recommended schemes and warn if not recommended
		u, perr := url.Parse(val)
		if perr == nil {
			scheme := strings.ToLower(u.Scheme)
			// Recommended schemes per the draft
			recommended := map[string]bool{"http": true, "https": true, "mailto": true, "tel": true}
			if !recommended[scheme] {
				// Moderate warning: syntactically allowed but not recommended
				res.Messages = append(res.Messages, fmt.Sprintf("furi uses non-recommended scheme %q; the draft RECOMMENDS only http, https, mailto and tel. Non-recommended schemes may be unsafe; do NOT auto-follow without user confirmation.", scheme))
				if scheme == "javascript" || scheme == "data" {
					res.Messages = append(res.Messages, "Note: this scheme can be dangerous (may execute code or embed data). Treat as potentially unsafe and require manual review before following.")
				}
			}
		}

		// As with ftxt, check that textual content is valid UTF-8 and free of disallowed control characters.
		warns, errs := checkUnicodeContent(val)
		for _, w := range warns {
			// for URIs, control characters are usually invalid; we treat errors strictly
			res.Messages = append(res.Messages, "Warning: "+w)
		}
		if len(errs) > 0 {
			for _, e := range errs {
				res.Messages = append(res.Messages, "Error: "+e)
			}
			// Let validateURI also run, but mark invalid because of disallowed characters
			res.Valid = false
			return res
		}

		// Now run existing validation (which provides additional syntax/semantic checks).
		if err := validateURI(val); err != nil {
			res.Messages = append(res.Messages, fmt.Sprintf("furi parsing error: %v", err))
			// Per spec: URIs MUST conform; but since version tag is present, processors MAY treat as for sale while warning.
			res.Valid = true
			res.Messages = append(res.Messages, "Because the version tag is present, processors SHOULD treat the domain as for sale even if the furi value is syntactically invalid. This tool marks the record as ACCEPTED with warnings.")
			return res
		}

		// If we reached here, the URI parsed and passed validateURI checks.
		res.Valid = true
		res.Messages = append(res.Messages, "furi content tag contains a syntactically valid URI/IRI. Do NOT auto-redirect users to this URI without prompting (security risk).")
		return res

	case "fval":
		if len(val) < 2 {
			res.Messages = append(res.Messages, "fval value too short (must be at least 2 characters: currency+amount).")
			res.Valid = false
			return res
		}
		if len([]byte(val)) > 239 {
			res.Messages = append(res.Messages, fmt.Sprintf("fval value byte length %d exceeds the draft's maximum of 239 characters for fval-value.", len([]byte(val))))
			res.Valid = false
			return res
		}
		if !fvalRe.MatchString(val) {
			res.Messages = append(res.Messages, "fval value does not conform to the required format: <CURRENCY><AMOUNT>, e.g. USD750 or BTC0.000010. Currency MUST be uppercase letters; amount MUST be digits with optional fractional part.")
			res.Valid = false
			return res
		}
		idx := firstDigitIndex(val)
		if idx <= 0 {
			res.Messages = append(res.Messages, "Unable to separate currency code and amount in fval value.")
			res.Valid = false
			return res
		}
		_ = val[:idx] // currency (not strictly enforcing ISO4217)
		res.Valid = true
		res.Messages = append(res.Messages, "fval content tag is syntactically acceptable. Note: prices are indicative only; verify with seller.")
		return res

	default:
		res.Messages = append(res.Messages, fmt.Sprintf("Unknown content tag: %q", foundTag))
		res.Valid = false
		return res
	}
}

// checkUnicodeContent checks that the given string is valid UTF-8 and
// flags the presence of control characters or non-characters according to the
// draft's recommendations.
//
// Returns two slices: warnings and errors. Warnings are moderate advisory notes
// (for example: presence of tab/CR/LF which are "best avoided"); errors are
// violations that should cause the record to be treated as invalid (e.g. other
// control characters, C1 controls, invalid UTF-8).
func checkUnicodeContent(s string) (warnings []string, errorsOut []string) {
	if !utf8.ValidString(s) {
		errorsOut = append(errorsOut, "content is not valid UTF-8; the draft RECOMMENDS UTF-8 encoding for text content")
		return
	}

	for i, r := range s {
		// C0 controls (U+0000..U+001F) and DEL (U+007F)
		if r <= 0x1F || r == 0x7F {
			// Exception per draft: U+0009 (TAB), U+000A (LF), U+000D (CR) are "best avoided" -> warn
			if r == 0x09 || r == 0x0A || r == 0x0D {
				warnings = append(warnings, fmt.Sprintf("contains control character U+%04X at byte index %d (TAB/CR/LF are allowed but RECOMMENDED to be avoided)", r, i))
			} else {
				errorsOut = append(errorsOut, fmt.Sprintf("contains disallowed control character U+%04X at byte index %d; other control characters are not permitted in content values", r, i))
			}
		}

		// C1 controls (U+0080..U+009F) are controls and should be considered invalid
		if r >= 0x80 && r <= 0x9F {
			errorsOut = append(errorsOut, fmt.Sprintf("contains C1 control U+%04X at byte index %d; C1 controls are not permitted", r, i))
		}

		// Non-characters: U+FDD0..U+FDEF and any codepoint where low 16 bits are 0xFFFE or 0xFFFF
		if (r >= 0xFDD0 && r <= 0xFDEF) || (r&0xFFFF == 0xFFFE) || (r&0xFFFF == 0xFFFF) {
			warnings = append(warnings, fmt.Sprintf("contains Unicode non-character U+%04X at index %d; non-characters are discouraged for interchange", r, i))
		}
	}

	return
}

// unescapePresentation decodes presentation-format escapes found in DNS zone file strings.
// It supports:
//   - \DDD where D are 1..3 decimal digits representing an octet value (0..255)
//   - backslash escaping of a single character: e.g. \"  \\  \;  etc.
//
// This matches common DNS zone-file presentation semantics (and handles the examples
// where TXT RDATA contains sequences like "\240\159\142\133" representing UTF-8 bytes).
func unescapePresentation(s string) ([]byte, error) {
	var out []byte
	i := 0
	for i < len(s) {
		c := s[i]
		if c != '\\' {
			out = append(out, c)
			i++
			continue
		}
		// backslash found
		i++
		if i >= len(s) {
			// stray backslash at end -> treat as literal backslash
			out = append(out, '\\')
			break
		}
		// If next is digit, parse up to 3 decimal digits
		if s[i] >= '0' && s[i] <= '9' {
			start := i
			end := i
			for end < len(s) && end-start < 3 && s[end] >= '0' && s[end] <= '9' {
				end++
			}
			numStr := s[start:end]
			var val int
			for _, ch := range numStr {
				val = val*10 + int(ch-'0')
			}
			if val < 0 || val > 255 {
				return out, fmt.Errorf("escaped decimal value out of range: %s", numStr)
			}
			out = append(out, byte(val))
			i = end
			continue
		}
		// Not digits: take the next character literally (as per presentation escaping)
		out = append(out, s[i])
		i++
	}
	return out, nil
}

func firstDigitIndex(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			return i
		}
	}
	return -1
}

// validateURI attempts to parse the value as a URI per RFC3986.
// net/url.Parse is used as a sanity check. A scheme is required.
// For http(s) URIs, checks for obviously invalid host characters (e.g., backslash).
func validateURI(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("url.Parse failed: %w", err)
	}
	if u.Scheme == "" {
		return errors.New("no URI scheme found (e.g., https, http, mailto, tel). A scheme is required for furi")
	}
	if strings.Contains(s, " ") {
		return fmt.Errorf("URI contains unencoded spaces")
	}
	if u.Scheme == "mailto" {
		if u.Opaque == "" && u.Path == "" {
			return errors.New("mailto: URI contains no recipient address")
		}
	}
	if u.Scheme == "tel" && u.Opaque == "" && u.Path == "" {
		return errors.New("tel: URI contains no telephone number")
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		host := u.Host
		// Strip optional port
		if strings.Contains(host, ":") {
			h, _, err := net.SplitHostPort(host)
			if err == nil {
				host = h
			}
		}
		if host == "" {
			return errors.New("http(s) URI has empty host")
		}
		if strings.ContainsAny(host, "\\\n\r\t\x00") {
			return errors.New("URI host contains invalid characters")
		}
	}
	return nil
}
