syntax: v1

gen:
    hooks:
        {{ if has "serverless" .Features }}before:
            - gorename {{ .Module }}/server {{ .Module }}/internal{{ end }}
        after:
            {{ if has "serverless" .Features }}- gorename {{ .Module }}/internal {{ .Module }}/server{{ end }}
            - jzero gen swagger

    split-api-types-dir: true