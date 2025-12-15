{{ if ne .Style "gozero" }}style: {{.Style}}

{{ end }}gen:
    hooks:
        after:
            - jzero gen swagger