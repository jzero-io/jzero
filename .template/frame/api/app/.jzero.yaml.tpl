syntax: v1

gen:
    hooks:
        {{ if or (has "serverless" .Features) (has "serverless_core" .Features) }}before:
            - gorename {{ .Module }}/server {{ .Module }}/internal{{ end }}
        after:
            {{ if or (has "serverless" .Features) (has "serverless_core" .Features) }}- gorename {{ .Module }}/internal {{ .Module }}/server{{ end }}
            - jzero gen swagger

    split-api-types-dir: true