{{if (VersionCompare .GoVersion ">=" "1.24")}}
hooks:
    before:
        - go install tool
{{end}}

gen:
    hooks:
        after:
            - jzero gen swagger

    style: {{.Style}}