package swaggerv2

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/eddieowens/opts"
	"github.com/zeromicro/go-zero/rest"

	"github.com/jzero-io/jzero/core/templatex"
)

const defaultSwaggerTemplate = `
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

type Swaggerv2Opts struct {
	SwaggerPath     string
	SwaggerHost     string
	SwaggerTemplate string
}

func WithSwaggerHost(swaggerHost string) opts.Opt[Swaggerv2Opts] {
	return func(config *Swaggerv2Opts) {
		config.SwaggerHost = swaggerHost
	}
}

func WithSwaggerPath(swaggerPath string) opts.Opt[Swaggerv2Opts] {
	return func(config *Swaggerv2Opts) {
		config.SwaggerPath = swaggerPath
	}
}

func WithSwaggerTemplate(swaggerTemplate string) opts.Opt[Swaggerv2Opts] {
	return func(config *Swaggerv2Opts) {
		config.SwaggerTemplate = swaggerTemplate
	}
}

func (opts Swaggerv2Opts) DefaultOptions() Swaggerv2Opts {
	return Swaggerv2Opts{
		SwaggerPath:     filepath.Join("desc", "swagger"),
		SwaggerHost:     "https://petstore.swagger.io",
		SwaggerTemplate: defaultSwaggerTemplate,
	}
}

func RegisterRoutes(server *rest.Server, op ...opts.Opt[Swaggerv2Opts]) {
	o := opts.DefaultApply(op...)

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger/:path",
		Handler: rawHandler(o),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger",
		Handler: uiHandler(o),
	})
}

func rawHandler(config Swaggerv2Opts) http.HandlerFunc {
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

func uiHandler(config Swaggerv2Opts) http.HandlerFunc {
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
		}, []byte(config.SwaggerTemplate))
		_, _ = rw.Write(uiHTML)
	}
}

func getSwaggerFiles(dir string) ([]string, error) {
	var files []string

	swaggerDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, swaggerFile := range swaggerDir {
		if swaggerFile.IsDir() {
			filenames, err := getSwaggerFiles(filepath.Join(dir, swaggerFile.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, filenames...)
		} else {
			if strings.HasSuffix(swaggerFile.Name(), ".json") {
				files = append(files, filepath.Join(swaggerFile.Name()))
			}
		}
	}
	// 保证如果有 swagger.json, 就放第一位
	for i, file := range files {
		if file == "swagger.json" {
			files[0], files[i] = files[i], files[0]
			break
		}
	}
	return files, nil
}
