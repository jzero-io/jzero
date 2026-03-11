---
title: mongodb guide
icon: devicon-plain:mongodb-wordmark
star: true
order: 5
---

## Introduction

jzero supports generating code to `internal/mongo` by specifying mongo type.

For easier usage, jzero automatically generates `internal/mongo/model.go` file to register all generated mongo models.

## Features

* Supports redis and custom cache

## Generate code

```yaml
gen:
    # specify mongo type
    mongo-type: ["user", "role", "menu"]
    # whether cache is needed
    mongo-cache: true
    # specify which types need cache
    mongo-cache-type:
      - user
```

```shell
jzero gen
```
