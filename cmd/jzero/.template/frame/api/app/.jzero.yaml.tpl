{{if (VersionCompare .GoVersion ">=" "1.24")}}
hooks:
    before:
        - go mod tidy
        - go install tool
{{end}}

gen:
    hooks:
        after:
            - jzero gen swagger
            - jzero format

    style: {{.Style}}