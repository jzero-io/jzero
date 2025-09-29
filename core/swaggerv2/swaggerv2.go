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
			{{range $k, $v := .SwaggerJsonsPath}}{url: "swagger?path={{ $v }}", name: "{{ $v }}"},
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
		Path:    "/swagger",
		Handler: SwaggerHandlerFunc(o),
	})
}

func SwaggerHandlerFunc(config Swaggerv2Opts) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		SwaggerHandler(config, rw, r)
	}
}

func SwaggerHandler(config Swaggerv2Opts, w http.ResponseWriter, r *http.Request) {
	if r.FormValue("path") != "" {
		file, err := os.ReadFile(filepath.Join(config.SwaggerPath, r.FormValue("path")))
		if err != nil {
			return
		}
		_, _ = w.Write(file)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, strings.TrimSuffix(r.RequestURI, "/"), 301)
		}

		swaggerJsonsPath, err := getSwaggerFiles(config.SwaggerPath)
		if err != nil {
			return
		}

		uiHTML, _ := templatex.ParseTemplate(map[string]any{
			"SwaggerHost":      config.SwaggerHost,
			"SwaggerJsonsPath": swaggerJsonsPath,
		}, []byte(config.SwaggerTemplate))
		_, _ = w.Write(uiHTML)
	}
}

func getSwaggerFiles(dir string) ([]string, error) {
	return getSwaggerFilesRecursive(dir, dir)
}

func getSwaggerFilesRecursive(rootDir, currentDir string) ([]string, error) {
	var files []string

	swaggerDir, err := os.ReadDir(currentDir)
	if err != nil {
		return nil, err
	}

	for _, swaggerFile := range swaggerDir {
		if swaggerFile.IsDir() {
			filenames, err := getSwaggerFilesRecursive(rootDir, filepath.Join(currentDir, swaggerFile.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, filenames...)
		} else {
			if strings.HasSuffix(swaggerFile.Name(), ".json") {
				// 计算相对于根目录的路径
				relPath, err := filepath.Rel(rootDir, filepath.Join(currentDir, swaggerFile.Name()))
				if err != nil {
					return nil, err
				}
				files = append(files, relPath)
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
