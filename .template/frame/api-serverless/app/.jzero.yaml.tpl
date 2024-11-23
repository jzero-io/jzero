syntax: v1

gen:
    hooks:
        before:
            - gorename {{ .Module }}/server {{ .Module }}/internal
        after:
            - gorename {{ .Module }}/internal {{ .Module }}/server
            - jzero gen swagger

    split-api-types-dir: true
    regen-api-handler: true