package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// DomainInfo struct contains all relevant information about the 'for-sale' status of a domain.
type DomainInfo struct {
	Domain     string
	ForSale    bool
	ValidTags  []string
	InvalidRaw []string
	ErrorMsg   string
}

var (
	validFVALChar = regexp.MustCompile(`^[A-Z]{3}[0-9.,]{1,236}$`)
)

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/check", checkHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("form").Parse(`
		<!DOCTYPE html>
		<html lang="en"><head><title>For Sale Check</title></head>
		<body>
		<h2>Check if a domain name is for sale</h2>
		<form action="/check" method="get">
			Domain Name: <input type="text" name="domain" value="example.nl">
			<input type="submit" value="Check">
		</form></body></html>
	`))
	tmpl.Execute(w, nil)
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	info := DomainInfo{Domain: domain}

	if domain == "" {
		info.ErrorMsg = "No domain name provided"
		renderResult(w, info)
		return
	}

	queryName := "_for-sale." + domain
	txts, err := net.LookupTXT(queryName)
	if err != nil {
		info.ErrorMsg = fmt.Sprintf("No _for-sale TXT records found: %v", err)
		renderResult(w, info)
		return
	}

	seen := map[string]bool{}
	foundValidVersionTag := false

	for _, txt := range txts {
		if seen[txt] {
			continue
		}
		seen[txt] = true

		if !strings.HasPrefix(txt, "v=FORSALE1;") {
			info.InvalidRaw = append(info.InvalidRaw, txt)
			continue
		}

		foundValidVersionTag = true
		info.ForSale = true

		content := strings.TrimPrefix(txt, "v=FORSALE1;")
		content = strings.TrimLeft(content, " ")

		if content == "" {
			info.ValidTags = append(info.ValidTags, "v=FORSALE1;")
			continue
		}

		parts := strings.SplitN(content, "=", 2)
		if len(parts) != 2 {
			info.InvalidRaw = append(info.InvalidRaw, txt)
			continue
		}

		tag, value := parts[0], parts[1]

		switch tag {
		case "fcod":
			if len(value) >= 1 && len(value) <= 239 {
				info.ValidTags = append(info.ValidTags, content)
			} else {
				info.InvalidRaw = append(info.InvalidRaw, txt)
			}
		case "ftxt":
			if len(value) >= 1 && len(value) <= 239 {
				info.ValidTags = append(info.ValidTags, content)
			} else {
				info.InvalidRaw = append(info.InvalidRaw, txt)
			}
		case "furi":
			u, err := url.Parse(value)
			if err == nil {
				scheme := strings.ToLower(u.Scheme)
				if scheme == "http" || scheme == "https" || scheme == "mailto" || scheme == "tel" {
					info.ValidTags = append(info.ValidTags, content)
					break
				}
			}
			info.InvalidRaw = append(info.InvalidRaw, txt)
		case "fval":
			if len(value) >= 4 && len(value) <= 239 && validFVALChar.MatchString(value) && strings.Count(value, ".") <= 1 {
				info.ValidTags = append(info.ValidTags, content)
			} else {
				info.InvalidRaw = append(info.InvalidRaw, txt)
			}
		default:
			info.InvalidRaw = append(info.InvalidRaw, txt)
		}
	}

	if !foundValidVersionTag {
		info.ForSale = false
	}

	renderResult(w, info)
}

func renderResult(w http.ResponseWriter, info DomainInfo) {
	funcMap := template.FuncMap{
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"stripPrefix": func(s, prefix string) string {
			return strings.TrimPrefix(s, prefix)
		},
		"htmlEscape": func(s string) string {
			return template.HTMLEscapeString(s)
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
	}

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
					<li><code>Text message: {{.}}</code></li>
				{{- else if hasPrefix . "fval=" -}}
					{{ $val := stripPrefix . "fval=" }}
					<li><code>Price: {{.}}</code></li>
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

	tmpl.Execute(w, info)
}
