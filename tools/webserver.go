package main

import (
	"fmt"          // For formatted I/O, such as printing to console and strings.
	"html/template" // For safe rendering of HTML templates, including XSS prevention.
	"log"          // For logging events (e.g., server start, fatal errors).
	"net"          // For network functionality, specifically DNS lookups (LookupTXT).
	"net/http"     // For setting up an HTTP server and handling requests.
	"net/url"      // For parsing and validating URLs (URIs).
	"regexp"       // For regular expressions, used for ftxt-character validation.
	"strings"      // For string manipulation, such as removing prefixes and splitting.
)

// DomainInfo struct contains all relevant information about the 'for-sale' status of a domain.
type DomainInfo struct {
	Domain     string   // The domain name being checked.
	ForSale    bool     // Indicates whether the domain is marked as 'for sale'.
	ValidTags  []string // List of syntactically valid and recognized tags (e.g., fcod=..., furi=...).
	InvalidRaw []string // List of records that had the 'v=FORSALE1;' tag but were otherwise invalid.
	ErrorMsg   string   // Error message if there was a problem with the DNS lookup or input.
}

var (
	// validFTXTChar is a regular expression to validate that ftxt-values only contain allowed characters.
	// RFC 3.1: "ftxt-char = %x20-21 / %x23-5B / %x5D-7E" (excluding " and \)
	validFTXTChar = regexp.MustCompile(`^[\x20-\x21\x23-\x5B\x5D-\x7E]+$`)
)

func main() {
	// Register the handler functions for the different URL paths.
	http.HandleFunc("/", formHandler)     // Handles the root URL (the input form).
	http.HandleFunc("/check", checkHandler) // Handles the '/check' URL (the DNS check).

	// Start the HTTP server on port 8080.
	log.Println("Server running on http://localhost:8080")
	// log.Fatal ensures the program exits if the server cannot start or a fatal error occurs.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// formHandler displays a simple HTML form for entering a domain name.
func formHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the HTML template for the form. template.Must panics if there's an error.
	tmpl := template.Must(template.New("form").Parse(`
		<!DOCTYPE html>
		<html lang="en"><head><title>For Sale Check</title></head>
		<body>
		<h2>Check if a domain name is for sale</h2>
		<form action="/check" method="get">
			Domain Name: <input type="text" name="domain">
			<input type="submit" value="Check">
		</form></body></html>
	`))
	// Execute the template and write the result to the http.ResponseWriter.
	tmpl.Execute(w, nil)
}

// checkHandler performs the DNS lookup and processes the '_for-sale' TXT records.
func checkHandler(w http.ResponseWriter, r *http.Request) {
	// Get the domain name from the URL parameters.
	domain := r.URL.Query().Get("domain")
	info := DomainInfo{Domain: domain} // Initialize the DomainInfo struct.

	// Check if a domain name was provided.
	if domain == "" {
		info.ErrorMsg = "No domain name provided"
		renderResult(w, info) // Render the result with an error message.
		return
	}

	// Construct the DNS query name according to the specification: "_for-sale.<domain>".
	queryName := "_for-sale." + domain
	txts, err := net.LookupTXT(queryName) // Perform the DNS TXT lookup.
	// TODO misses the _for-sale IN TXT "v=FORSALE1;" "ftxt=foo" "bar" "invalid" case, because it concatenates implicitly
	if err != nil {
		// RFC 3.1: "If no TXT records at a leaf node contain a valid version tag,
		// processors MUST consider the node name invalid and discard it."
		// This means that if no records are found (NoAnswer) or the domain does not exist (NXDOMAIN),
		// the domain is not considered 'for sale' via this method.
		info.ErrorMsg = fmt.Sprintf("No _for-sale TXT records found: %v", err)
		renderResult(w, info)
		return
	}

	// Use a map to keep track of already processed (duplicate) TXT records.
	// TODO Does this make sense?
	seen := map[string]bool{}
	// foundValidVersionTag keeps track of whether at least one TXT record with a valid 'v=FORSALE1;' tag was found.
	// This is crucial for the robustness principle of the RFC.
	foundValidVersionTag := false

	// Loop through all found TXT records.
	for _, txt := range txts {
		// Skip duplicates to avoid redundant processing and output.
		if seen[txt] {
			continue
		}
		seen[txt] = true

		// RFC 3.1: "TXT records in the same RRset, but without a version tag,
		// MUST NOT be interpreted or processed as a valid '_for-sale' indicator."
		if !strings.HasPrefix(txt, "v=FORSALE1;") {
			// These records are added to InvalidRaw as "noise", but do not lead to the 'for sale' status.
			info.InvalidRaw = append(info.InvalidRaw, txt)
			continue
		}

		// If we reach this point, the record has the mandatory 'v=FORSALE1;' version tag.
		// RFC 3.1: "If the version tag itself is valid, processors MAY assume that the domain is for sale."
		// Therefore, we set info.ForSale to true directly. This remains true even if further content turns out to be invalid.
		foundValidVersionTag = true
		info.ForSale = true

		// Remove the version tag from the beginning of the string.
		content := strings.TrimPrefix(txt, "v=FORSALE1;")
		// Remove any leading spaces after the version tag.
		// RFC 5: "This also applies to space characters (%x20) immediately following the version tag."
		content = strings.TrimLeft(content, " ")

		// RFC 3.1: "In the absence of a tag-value pair, processors MAY assume that the domain is for sale."
		// If, after removing the version tag and spaces, the content is empty,
		// it means the record only contained "v=FORSALE1;". This is valid.
		if content == "" {
			info.ValidTags = append(info.ValidTags, "v=FORSALE1;") // Add the bare version tag as a valid indicator.
			continue                                             // Continue with the next record.
		}

		// Try to split the rest of the content into a 'tag' and a 'value' at the first '=' sign.
		// RFC 3.1: "only one tag-value pair per record".
		parts := strings.SplitN(content, "=", 2)
		if len(parts) != 2 {
			// RFC 3.1: "If a tag-value pair is present but invalid, this constitutes a syntax
			// error and SHOULD be treated as if it were absent. In such cases, if the
			// version tag itself is valid, processors MAY assume that the domain is for sale."
			// The record is marked as invalid, but because foundValidVersionTag is already true,
			// info.ForSale remains true. We immediately proceed to the next record,
			// as further parsing of this invalid content is pointless.
			info.InvalidRaw = append(info.InvalidRaw, txt)
			continue
		}

		tag, value := parts[0], parts[1] // Split into tag and value.

		// Evaluate the tag and validate the corresponding value.
		switch tag {
		case "fcod": // 'For Sale Code' tag.
			// RFC 3.1: "fcod-value = 1*239OCTET" (minimum 1, maximum 239 octets).
			if len(value) >= 1 && len(value) <= 239 {
				info.ValidTags = append(info.ValidTags, content) // Add the entire tag=value string to valid tags.
			} else {
				info.InvalidRaw = append(info.InvalidRaw, txt) // Invalid length.
			}
		case "ftxt": // 'Free Text' tag.
			// RFC 3.1: "ftxt-value = 1*239ftxt-char" and "ftxt-char = %x20-21 / %x23-5B / %x5D-7E"
			// (excluding double quotes " and backslash \ to prevent escape issues).
			if len(value) >= 1 && len(value) <= 239 && validFTXTChar.MatchString(value) {
				info.ValidTags = append(info.ValidTags, content)
			} else {
				info.InvalidRaw = append(info.InvalidRaw, txt) // Invalid length or incorrect characters.
			}
		case "furi": // 'For Sale URI' tag.
			u, err := url.Parse(value) // Try to parse the value as a URI.
			// RFC 3.1: "Only http, https, mailto and tel schemes" are recommended.
			if err == nil { // Only if the URI is syntactically valid.
				scheme := strings.ToLower(u.Scheme) // Get the scheme and convert to lowercase.
				if scheme == "http" || scheme == "https" || scheme == "mailto" || scheme == "tel" {
					info.ValidTags = append(info.ValidTags, content)
					break // Exit the switch statement, as the URI is correct.
				}
			}
			info.InvalidRaw = append(info.InvalidRaw, txt) // Invalid URI or disallowed scheme.
		default:
			// This catches unknown tags (e.g., "v=FORSALE1;foo=bar").
			// RFC 3.1: "If a tag-value pair is present but invalid, this constitutes a syntax
			// error and SHOULD be treated as if it were absent."
			// The record is marked as invalid, but 'info.ForSale' remains 'true' (thanks to the version tag).
			info.InvalidRaw = append(info.InvalidRaw, txt)
		}
	}

	// After iterating through all TXT records: if no record with a valid 'v=FORSALE1;' tag was found,
	// then the domain is not 'for sale' according to this specification.
	// This is necessary for cases where TXT records exist, but none meet the version requirement.
	if !foundValidVersionTag {
		info.ForSale = false
	}

	renderResult(w, info) // Render the final result to the client.
}

// renderResult generates the HTML output to display the results of the DNS check.
func renderResult(w http.ResponseWriter, info DomainInfo) {
	// Template.FuncMap defines custom functions that can be used within the HTML template.
	funcMap := template.FuncMap{
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"stripPrefix": func(s, prefix string) string {
			return strings.TrimPrefix(s, prefix)
		},
		"htmlEscape": func(s string) string {
			// template.HTMLEscapeString escapes HTML special characters in a string.
			// This is essential for XSS prevention when text is placed directly into the HTML body.
			return template.HTMLEscapeString(s)
		},
		"safeURL": func(s string) template.URL {
			// template.URL tells the template engine that the string is a safe URL
			// and does not need further escaping for use in URL attributes (like href).
			// This is necessary for schemes like 'tel:' or 'mailto:' that would otherwise be considered unsafe.
			return template.URL(s)
		},
	}

	// Parse the HTML template for the results page.
	tmpl := template.Must(template.New("result").Funcs(funcMap).Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head><title>For Sale Check</title></head>
		<body>
		<h2>Result for {{.Domain}}</h2>
		{{if .ErrorMsg}}
			<p style="color: red;">Error: {{.ErrorMsg}}</p>
		{{else if .ForSale}}
			<p style="color: green;">✅ Domain appears to be for sale based on found records.</p>
			<ul>
			{{range .ValidTags}}
				{{- if hasPrefix . "furi=" -}}
					{{ $uri := stripPrefix . "furi=" }}
					<li><a href="{{safeURL $uri}}" target="_blank" rel="noopener noreferrer">{{htmlEscape $uri}}</a> - click at own risk!</li>
				{{- else if hasPrefix . "fcod=" -}}
					<li><code style="color: #888;">{{.}}</code></li>
                                {{- else if hasPrefix . "ftxt=" -}}
                                        {{ $txt := stripPrefix . "ftxt=" }}
                                        <li><code>Text message: {{htmlEscape $txt}}</code></li>
				{{- else -}}
					<li><code>{{.}}</code></li>
				{{- end }}
			{{end}}
			</ul>
			{{if .InvalidRaw}}
				<p style="color: orange;">⚠️ Some records were syntactically invalid.</p>
				<ul>{{range .InvalidRaw}}<li><code>{{.}}</code></li>{{end}}</ul>
			{{end}}
		{{else}}
			<p>❌ No valid indications found that the domain is for sale.</p>
		{{end}}
		<a href="/">Back</a>
		</body></html>
	`))

	// Execute the template with the DomainInfo data.
	tmpl.Execute(w, info)
}
