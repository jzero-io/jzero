---
title: Template guide
icon: /icons/vscode-icons-folder-type-template.svg
star: true
order: 5.3
---

## Template initialization

Initialize jzero embedded templates or remote repository templates to local disk.

```shell
# Initialize jzero embedded templates to $HOME/.jzero/templates/$Version, can modify templates then create new projects
jzero template init
# Or initialize templates to current project's .template, jzero gen will prioritize reading current project's .template as template home
jzero template init --output .template
# Initialize remote repository templates to $HOME/.jzero/templates/remote, such as gateway
jzero template init --branch gateway

# If still need to extend go-zero's template
goctl template init --home .template/go-zero
```

## Initialize project with custom template

* Specify remote repository template

```shell
jzero new project_name --remote repo_to_your_templates --branch template_branch
# Get remote template from cache
jzero new project_name --remote repo_to_your_templates --branch template_branch --cache
```

* Use local template

```shell
jzero new project_name --local template_name
```

* Use path template

```shell
jzero new project_name --home path_to_template
```

## Template rendering and variables

`jzero new` renders both template content and template paths when creating a project:

* `.tpl` file contents are rendered with Go `text/template`
* File names and directory names are also rendered, so variables like `{{ .APP }}` and `{{ .Module }}` can be used in paths
* If a file ends with `.tpl.tpl`, only one `.tpl` suffix is removed and the file content is not rendered again, which is useful when you want to keep template source text

For example, this template path:

```text
internal/{{ .APP | lower }}/{{ FormatStyle .Style "service_context.go.tpl" }}
```

will be rendered into a real directory and file name during project creation.

### Built-in variables

When `jzero new` runs, jzero injects the following built-in variables into the template:

| Variable | Type | Description |
| --- | --- | --- |
| `APP` | `string` | Project name from `jzero new <name>` or `--name` |
| `Module` | `string` | Go module name from `--module`; defaults to the project name when omitted |
| `GoVersion` | `string` | Current Go version |
| `GoArch` | `string` | Current architecture such as `amd64` or `arm64` |
| `DirName` | `string` | Output directory name |
| `Style` | `string` | File naming style, default is `gozero` |
| `Features` | `[]string` | Feature list passed by `jzero new --features` |
| `Serverless` | `bool` | Whether the project is created in serverless mode |

Example:

```text
module {{ .Module }}

{{ if has "model" .Features }}
// enable model feature
{{ end }}

{{ if .Serverless }}
// serverless mode
{{ end }}
```

:::tip
`jzero template build` automatically rewrites the `go.mod` module and Go import paths that point to the current project into `{{ .Module }}`. That means templates produced by `jzero template build` can directly reuse the `Module` variable.
:::

### Built-in functions

Templates are rendered with Go `text/template`. In addition to built-in functions such as `and`, `or`, `not`, and `index`, you can also use many common functions from [sprig](https://masterminds.github.io/sprig/), such as `lower`, `upper`, `default`, `has`, and `dict`. jzero also registers these extra functions:

| Function | Description |
| --- | --- |
| `FirstUpper(s)` | Uppercase the first letter |
| `FirstLower(s)` | Lowercase the first letter |
| `ToCamel(s)` | Convert `foo-bar`, `foo_bar`, or `foo/bar` into camel case |
| `FormatStyle(style, name)` | Convert a file name using the selected `--style` |
| `VersionCompare(v1, op, v2)` | Compare versions with `>`, `<`, `>=`, `<=` |

Example:

```text
{{ .APP | ToCamel | FirstUpper }}
{{ FormatStyle .Style "service_context.go.tpl" }}
{{ if (VersionCompare .GoVersion ">=" "1.24") }}toolchain go1.24.0{{ end }}
```

### Inject custom template variables

Use the global `--register-tpl-val key=value` flag to inject extra template variables. Injected values are merged into the current template data, so they can be used in both template content and template paths.

```shell
jzero new myapi --local myapi \
  --register-tpl-val company=acme \
  --register-tpl-val owner=platform
```

You can access them directly in templates:

```text
# {{ .APP }}
Company: {{ .company }}
Owner: {{ .owner }}
```

They can also be used in paths:

```text
internal/{{ .company }}/banner.txt.tpl
```

If you want to reuse these variables across commands, put them in `.jzero.yaml`:

```yaml
register-tpl-val:
  - company=acme
  - owner=platform
```

Notes:

* If a custom variable has the same name as a built-in one, the custom value overrides it
* Values are currently parsed as `key=value`, so it is best not to include `=` in the value
* `--register-tpl-val` is a global flag. It is not limited to `jzero new`; other jzero commands that render templates also merge these values, but each command may provide different built-in variables

## Practice: Build your own template

:::tip Can convert any current project to jzero template, this is very cool!
:::

```bash
# Add a new api project
jzero new simpleapi
# Enter project
cd simpleapi
# Add a new api
jzero add api helloworld
# Generate code
jzero gen

# Build current project as template, save to $HOME/.jzero/templates/local/myapi
jzero template build --name myapi

# Now you can use your own template, you'll find the generated project automatically has helloworld api
jzero new mysimpleapi --local myapi

# But you find this template only allows local use, for universal effect
# You can create a templates repository in remote repository like github (assume https://github.com/jzero-io/templates)
# Then put content from $HOME/.jzero/templates/local/myapi into repository, and upload to myapi branch
jzero new project_name --remote https://github.com/jzero-io/templates --branch myapi
```

Template structure:

```bash
$ tree ~/.jzero/templates/local/myapi
└── app
    ├── Dockerfile.tpl
    ├── README.md.tpl
    ├── cmd
    │   ├── root.go.tpl
    │   ├── server.go.tpl
    │   └── version.go.tpl
    ├── desc
    │   ├── api
    │   │   ├── helloworld.api.tpl
    │   │   └── version.api.tpl
    │   └── swagger
    │       ├── helloworld.swagger.json.tpl
    │       ├── swagger.json.tpl
    │       └── version.swagger.json.tpl
    ├── etc
    │   └── etc.yaml.tpl
    ├── go.mod.tpl
    ├── internal
    │   ├── config
    │   │   └── config.go.tpl
    │   ├── custom
    │   │   └── custom.go.tpl
    │   ├── handler
    │   │   ├── helloworld
    │   │   │   └── helloworld_compact.go.tpl
    │   │   ├── routes.go.tpl
    │   │   └── version
    │   │       └── version.go.tpl
    │   ├── logic
    │   │   ├── helloworld
    │   │   │   └── create.go.tpl
    │   │   └── version
    │   │       └── version.go.tpl
    │   ├── middleware
    │   │   ├── middleware.go.tpl
    │   │   ├── response.go.tpl
    │   │   └── validator.go.tpl
    │   ├── svc
    │   │   ├── config.go.tpl
    │   │   ├── middleware.go.tpl
    │   │   └── servicecontext.go.tpl
    │   └── types
    │       ├── helloworld
    │       │   └── types.go.tpl
    │       ├── types.go.tpl
    │       └── version
    │           └── types.go.tpl
    ├── main.go.tpl
    └── plugins
        └── plugins.go.tpl
```
