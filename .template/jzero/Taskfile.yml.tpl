version: "3"

tasks:
  run:
    cmds:
      - go run main.go daemon
  build:
    cmds:
      - goreleaser build --snapshot --single-target --clean
    silent: true
  build:all:
    cmds:
      - goreleaser build --snapshot --clean
    silent: true
  build:amd64:
    cmds:
      - GOOS=linux GOARCH=amd64 goreleaser build --snapshot --single-target --clean
    silent: true
  build:arm64:
    cmds:
      - GOOS=linux GOARCH=arm64 goreleaser build --snapshot --single-target --clean
    silent: true