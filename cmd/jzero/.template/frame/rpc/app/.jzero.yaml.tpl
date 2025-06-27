{{if (VersionCompare .GoVersion ">=" "1.24")}}
hooks:
    before:
        - go mod tidy
        - go install tool
{{end}}

gen:
    style: {{.Style}}