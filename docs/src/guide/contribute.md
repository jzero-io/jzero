---
title: 贡献指南
icon: ooui:user-contributions-ltr
star: true
order: 30
category: contribute
tag:
  - contribute
---

欢迎参与 jzero 的开发以及维护, 这是一件非常有意义的事情, 让我们一起让 jzero 变得更好.

## 步骤

### 1. fork jzero

https://github.com/jzero-io/jzero/fork

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
go install github.com/fsgo/go_fmt/cmd/gorgeous@latest
gorgeous ./...
```

### 4. lint codes

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

1. fork jzero 并 clone jzero 到本地后

```shell
cd jzero
go install
```

2. new project with branch, e.g. `api`

```shell
jzero new your_project --branch api
```

3. run jzero gen with debug mode

```shell
jzero gen --debug --debug-sleep-time 15
```

4. attach jzero process

推荐采用 goland, 使用 attach 到 jzero 的进程中, 即可 debug, 如下所示:

<video width="720" height="450" controls>
  <source src="https://oss.jaronnie.com/iShot_2024-09-20_09.22.54.mp4" type="video/mp4">
</video>







