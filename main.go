package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// Package represents a Go package mapping
type Package struct {
	Path string // URL path like "/go-utils"
	Repo string
	VCS  string // "git", "hg", etc.
}

// Package mappings
var packages = []Package{
	{
		Path: "/go-utils",
		Repo: "https://github.com/aykhans/go-utils",
		VCS:  "git",
	},
	{
		Path: "/sarin",
		Repo: "https://github.com/aykhans/sarin",
		VCS:  "git",
	},
}

type HomeData struct {
	Domain   string
	Packages []Package
	Count    int
}

type VanityData struct {
	Domain string
	Path   string
	Repo   string
	VCS    string
}

type RedirectData struct {
	Domain string
	Path   string
}

var hostHeader = getEnvOrDefault("HOST_HEADER", "X-Forwarded-Host")
var hostOverride = os.Getenv("HOST")

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getHost(r *http.Request) string {
	if hostOverride != "" {
		return hostOverride
	}
	if host := r.Header.Get(hostHeader); host != "" {
		return host
	}
	return r.Host
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := HomeData{
		Domain:   getHost(r),
		Packages: packages,
		Count:    len(packages),
	}

	tmpl, err := template.New("home").Parse(homeTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func packageHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Find package by path
	var pkg *Package
	for i := range packages {
		if packages[i].Path == path {
			pkg = &packages[i]
			break
		}
	}

	if pkg == nil {
		http.NotFound(w, r)
		return
	}

	// Check if this is a go-get request
	if r.URL.Query().Get("go-get") == "1" {
		// Respond with meta tags for go get
		data := VanityData{
			Domain: getHost(r),
			Path:   pkg.Path,
			Repo:   pkg.Repo,
			VCS:    pkg.VCS,
		}

		tmpl, err := template.New("vanity").Parse(vanityTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Template execution error: %v", err)
		}
	} else {
		// Regular browser request - redirect to pkg.go.dev
		data := RedirectData{
			Domain: getHost(r),
			Path:   pkg.Path,
		}

		tmpl, err := template.New("redirect").Parse(redirectTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Template execution error: %v", err)
		}
	}
}

func main() {
	// Home page
	http.HandleFunc("/{$}", homeHandler)

	// All package paths
	http.HandleFunc("/", packageHandler)

	port := getEnvOrDefault("PORT", "8421")
	log.Printf("Server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

const vanityTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta name="go-import" content="{{.Domain}}{{.Path}} {{.VCS}} {{.Repo}}">
    <meta name="go-source" content="{{.Domain}}{{.Path}} {{.Repo}} {{.Repo}}/tree/master{/dir} {{.Repo}}/blob/master{/dir}/{file}#L{line}">
</head>
<body>
    go get {{.Domain}}{{.Path}}
</body>
</html>`

const redirectTemplate = `<!DOCTYPE html>
<html>
<head>
    <link rel="icon" href="https://pkg.go.dev/static/shared/icon/favicon.ico">
    <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/{{.Domain}}{{.Path}}">
</head>
<body>
    Redirecting to <a href="https://pkg.go.dev/{{.Domain}}{{.Path}}">pkg.go.dev/{{.Domain}}{{.Path}}</a>...
</body>
</html>`

const homeTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <link rel="icon" href="https://pkg.go.dev/static/shared/icon/favicon.ico">
    <title>Go Packages - {{.Domain}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            line-height: 1.6;
            color: #333;
        }
        h1 {
            color: #00ADD8;
            border-bottom: 2px solid #00ADD8;
            padding-bottom: 10px;
        }
        .package {
            background: #f5f5f5;
            border-left: 4px solid #00ADD8;
            padding: 15px;
            margin: 15px 0;
            border-radius: 4px;
        }
        .package-name {
            font-family: monospace;
            font-size: 1.1em;
            color: #00ADD8;
            font-weight: bold;
            text-decoration: none;
        }
        .package-name:hover {
            text-decoration: underline;
        }
        .external-icon {
            margin-left: 5px;
            vertical-align: middle;
        }
        .package-repo {
            margin-top: 5px;
            font-size: 0.9em;
        }
        .package-repo a {
            color: #666;
            text-decoration: none;
        }
        .package-repo a:hover {
            text-decoration: underline;
        }
        .install-cmd-wrapper {
            position: relative;
            margin-top: 8px;
            display: flex;
            align-items: center;
        }
        .install-cmd {
            background: #2d2d2d;
            color: #f8f8f2;
            padding: 10px;
            padding-right: 45px;
            border-radius: 4px;
            font-family: monospace;
            overflow-x: auto;
            flex: 1;
        }
        .copy-btn {
            position: absolute;
            right: 8px;
            background: #444;
            border: none;
            border-radius: 4px;
            padding: 6px 10px;
            cursor: pointer;
            color: #f8f8f2;
            font-size: 14px;
            transition: background 0.2s;
            display: flex;
            align-items: center;
            justify-content: center;
            height: 28px;
        }
        .copy-btn:hover {
            background: #555;
        }
        .copy-btn:active {
            background: #00ADD8;
        }
        .copy-btn.copied {
            background: #00ADD8;
        }
        .footer {
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #ddd;
            text-align: center;
            color: #666;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <h1>Go Packages</h1>

    {{range $index, $pkg := .Packages}}
    <div class="package">
        <a class="package-name" href="https://pkg.go.dev/{{$.Domain}}{{$pkg.Path}}" target="_blank">{{$.Domain}}{{$pkg.Path}}<svg class="external-icon" width="14" height="14" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M6 3H3a1 1 0 0 0-1 1v9a1 1 0 0 0 1 1h9a1 1 0 0 0 1-1v-3M9 2h5m0 0v5m0-5L7 9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></a>
        <div class="package-repo">
            Source: <a href="{{$pkg.Repo}}" target="_blank">{{$pkg.Repo}}<svg class="external-icon" width="12" height="12" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M6 3H3a1 1 0 0 0-1 1v9a1 1 0 0 0 1 1h9a1 1 0 0 0 1-1v-3M9 2h5m0 0v5m0-5L7 9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/></svg></a>
        </div>
        <div class="install-cmd-wrapper">
            <div class="install-cmd" id="cmd-{{$index}}">go get {{$.Domain}}{{$pkg.Path}}</div>
            <button class="copy-btn" onclick="copyToClipboard('cmd-{{$index}}', this)">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M13.5 5.5v7a1.5 1.5 0 0 1-1.5 1.5H5a1.5 1.5 0 0 1-1.5-1.5v-7A1.5 1.5 0 0 1 5 4h7a1.5 1.5 0 0 1 1.5 1.5z" stroke="currentColor" stroke-width="1.5" fill="none"/>
                    <path d="M5 4V3.5A1.5 1.5 0 0 1 6.5 2h7A1.5 1.5 0 0 1 15 3.5v7a1.5 1.5 0 0 1-1.5 1.5H13" stroke="currentColor" stroke-width="1.5" fill="none"/>
                </svg>
            </button>
        </div>
    </div>
    {{end}}

    <div class="footer">
        <p>Total packages: {{.Count}}</p>
    </div>

    <script>
        function copyToClipboard(elementId, button) {
            const element = document.getElementById(elementId);
            const text = element.textContent;

            navigator.clipboard.writeText(text).then(() => {
                // Visual feedback
                button.classList.add('copied');
                const originalHTML = button.innerHTML;
                button.innerHTML = '<svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M13 4L6 11L3 8" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>';

                setTimeout(() => {
                    button.classList.remove('copied');
                    button.innerHTML = originalHTML;
                }, 2000);
            }).catch(err => {
                console.error('Failed to copy:', err);
            });
        }
    </script>
</body>
</html>`
