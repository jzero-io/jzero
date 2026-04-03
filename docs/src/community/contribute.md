---
title: Welcome to Contribute👏
icon: /icons/ooui-user-contributions-ltr.svg
star: true
order: 30
---

Welcome to participate in jzero's development and maintenance. This is a very meaningful thing. Let's make jzero better together.

## Steps

### 1. fork jzero

[Click here to fork](https://github.com/jzero-io/jzero/fork)

### 2. clone

```shell
git clone https://github.com/your_username/jzero
```

### 3. checkout branch

```shell
cd jzero

git checkout -b feat/patch-1
```

### 4. format the code what you changes

```shell
jzero format
```

### 5. lint codes

```shell
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run --fix
```

### 6. push

```shell
git add .
git commit -m "feat(xx): custom message"
git push
```

### 7. pull request

Create your pull request!!!

## debug jzero

1. After forking jzero and cloning jzero locally

```shell
cd jzero
go install
```

2. new project with frame, e.g. `api`

```shell
jzero new your_project --frame api
```

3. run jzero gen with debug mode

```shell
jzero gen --debug --debug-sleep-time 15
```

4. attach jzero process

It's recommended to use goland, attach to jzero's process for debugging, as shown below:

<video width="720" height="450" controls>
  <source src="https://oss.jaronnie.com/iShot_2024-09-20_09.22.54.mp4" type="video/mp4">
</video>
