package swagger

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path"
)

type Opts struct {
	// BasePath for the UI path, defaults to: /
	BasePath string
	// Path combines with BasePath for the full UI path, defaults to: docs
	Path string
	// SpecURL the url to find the spec for
	SpecURL string
	// StylesURL url for swagger-ui.css
	StylesURL string
	// BundleURL url for swagger-ui-bundle.js
	BundleURL string
	// StandalonePresetURL url for swagger-ui-standalone-preset.js
	StandalonePresetURL string
	// Title for the documentation site, default to: API documentation
	Title string
}

// EnsureDefaults in case some options are missing
func (r *Opts) EnsureDefaults() {
	if r.BasePath == "" {
		r.BasePath = "/"
	}
	if r.Path == "" {
		r.Path = "swagger"
	}
	if r.SpecURL == "" {
		r.SpecURL = "/swagger.json"
	}
	if r.StylesURL == "" {
		r.StylesURL = swaggerStylesLatest
	}
	if r.BundleURL == "" {
		r.BundleURL = swaggerBundleLatest
	}
	if r.StandalonePresetURL == "" {
		r.StandalonePresetURL = swaggerStandalonePresetLatest
	}
	if r.Title == "" {
		r.Title = "API documentation"
	}
}

func Middleware(opts *Opts, next http.Handler) http.Handler {
	opts.EnsureDefaults()

	pth := path.Join(opts.BasePath, opts.Path)
	tmpl := template.Must(template.New("swagger").Parse(swaggerTemplate))

	buf := bytes.NewBuffer(nil)
	_ = tmpl.Execute(buf, opts)
	b := buf.Bytes()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == pth {
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
			rw.WriteHeader(http.StatusOK)

			_, _ = rw.Write(b)
			return
		}

		if next == nil {
			rw.Header().Set("Content-Type", "text/plain")
			rw.WriteHeader(http.StatusNotFound)
			_, _ = rw.Write([]byte(fmt.Sprintf("%q not found", pth)))
			return
		}
		next.ServeHTTP(rw, r)
	})
}

const (
	swaggerStylesLatest           = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@latest/swagger-ui.css"
	swaggerBundleLatest           = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@latest/swagger-ui-bundle.js"
	swaggerStandalonePresetLatest = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@latest/swagger-ui-standalone-preset.js"
	swaggerTemplate               = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
    <link rel="stylesheet" type="text/css" href="{{ .StylesURL }}">
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="{{ .BundleURL }}"></script>
<script src="{{ .StandalonePresetURL }}"></script>
<script>
    window.onload = function () {
        window.ui = SwaggerUIBundle({
            url: "{{ .SpecURL }}",
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            plugins: [
                SwaggerUIBundle.plugins.DownloadUrl
            ],
            layout: "StandaloneLayout"
        })
    }
</script>
</body>
</html>
`
)
