package swaggerv2

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/rest"

	"github.com/jzero-io/jzero/core/templatex"
)

type Opts func(*swaggerConfig)

// SwaggerOpts configures the Doc middlewares.
type swaggerConfig struct {
	// SwaggerPath the path to find the spec for
	SwaggerPath string

	// SwaggerHost for the js that generates the swagger ui site, defaults to: http://petstore.swagger.io/
	SwaggerHost string
}

func RegisterRoutes(server *rest.Server, opts ...Opts) {
	config := &swaggerConfig{
		SwaggerPath: filepath.Join("desc", "swagger"),
		SwaggerHost: "https://petstore.swagger.io",
	}
	for _, opt := range opts {
		opt(config)
	}

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger/:path",
		Handler: rawHandler(config),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger",
		Handler: uiHandler(config),
	})
}

func rawHandler(config *swaggerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var swaggerPath string
		err := filepath.Walk(config.SwaggerPath, func(path string, info os.FileInfo, err error) error {
			if info.Name() == filepath.Base(r.URL.Path) {
				swaggerPath = path
			}
			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		file, err := os.ReadFile(swaggerPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(file)
	}
}

func uiHandler(config *swaggerConfig) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")

		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(rw, r, strings.TrimSuffix(r.RequestURI, "/"), 301)
		}

		swaggerJsonsPath, err := getSwaggerFiles(config.SwaggerPath)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		uiHTML, _ := templatex.ParseTemplate(map[string]any{
			"SwaggerHost":      config.SwaggerHost,
			"SwaggerJsonsPath": swaggerJsonsPath,
		}, []byte(swaggerTemplateV2))
		_, _ = rw.Write(uiHTML)
	}
}

func getSwaggerFiles(dir string) ([]string, error) {
	var files []string

	protoDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, protoFile := range protoDir {
		if protoFile.IsDir() {
			filenames, err := getSwaggerFiles(filepath.Join(dir, protoFile.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, filenames...)
		} else {
			if strings.HasSuffix(protoFile.Name(), ".json") {
				files = append(files, filepath.Join(protoFile.Name()))
			}
		}
	}
	return files, nil
}

const swaggerTemplateV2 = `
	<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>API documentation</title>
    <link rel="stylesheet" type="text/css" href="{{ .SwaggerHost }}/swagger-ui.css" >
    <link rel="icon" type="image/png" href="{{ .SwaggerHost }}/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="{{ .SwaggerHost }}/favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="{{ .SwaggerHost }}/swagger-ui-bundle.js"> </script>
    <script src="{{ .SwaggerHost }}/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        "dom_id": "#swagger-ui",
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
		validatorUrl: null,
        urls: [
			{{range $k, $v := .SwaggerJsonsPath}}{url: "swagger/{{ $v }}", name: "{{ $v }}"},
			{{end}}
		]
      })

      // End Swagger UI call region
      window.ui = ui
    }
  </script>
  </body>
</html>`
