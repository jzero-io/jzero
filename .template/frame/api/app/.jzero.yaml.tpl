syntax: v1

gen:
    hooks:
        after:
            - jzero gen swagger

    split-api-types-dir: true