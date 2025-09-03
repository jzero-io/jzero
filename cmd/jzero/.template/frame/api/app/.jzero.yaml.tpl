gen:
    hooks:
        after:
            - jzero gen swagger
            - jzero format

    style: {{.Style}}