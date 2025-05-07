syntax: v1

gen:
    hooks:
        before:
            - goctl env -w GOCTL_EXPERIMENTAL=off
        after:
            - jzero gen swagger