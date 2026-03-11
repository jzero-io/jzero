---
title: Template guide
icon: vscode-icons:folder-type-template
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
