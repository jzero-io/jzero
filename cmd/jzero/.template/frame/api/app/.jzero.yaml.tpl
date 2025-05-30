syntax: v1

gen:
    hooks:
        after:
            - jzero gen swagger
            - jzero format