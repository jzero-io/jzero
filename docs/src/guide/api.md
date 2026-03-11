---
title: Comprehensive API Tutorial
icon: eos-icons:api
star: true
order: 0.2
---

## Overview

api is go-zero's self-developed domain-specific language (hereinafter referred to as api language or api descriptor language), aimed at implementing a user-friendly basic descriptor language as the most basic descriptor language for generating HTTP services.

jzero has extended api syntax, supporting the following features:
* `go_package`: Generates go types in defined packages, allowing different api files to have same-named type definitions, consistent with proto's `go_package`
* `compact_handler`: Generates handlers of the same route group in one file, reducing file count, consistent with proto's server module

## api definition

```api
info (
    // Define go_package: folder location for generated go types
    go_package: "user"
)

type User {
    id int `json:"id"`
    username string `json:"username"`
}

type PageRequest {
    page int `form:"page"`
    size int `form:"size"`
    username string `form:"username,optional"` // filter parameter, optional
}

type PageResponse {
    total uint64 `json:"total"` // total
    list  []User `json:"list"`  // paginated data
}

type UpdateRequest {
    id int `path:"id"`
    username string `json:"username"`
}

type UpdateResponse {}

@server (
    prefix:          /api/user // route prefix
    group:           user      // generated handler/logic folder location
    jwt:             JwtAuth   // whether to enable jwt authentication
    middleware:      AuthX     // middleware for this route group
    compact_handler: true      // whether to merge this group's handlers into one file, default each route has handler file
)
service simpleapi {
    @doc "User pagination"
    @handler Page
    get /page (PageRequest) returns (PageResponse)

    @doc "Update user"
    @handler Update
    post /update (UpdateRequest) returns (UpdateResponse)
}
```

Corresponding curl commands:

```shell
# User pagination endpoint
curl -X GET "http://localhost:8080/api/user/page?page=1&size=10&username=test" \
  -H "Authorization: Bearer <your-jwt-token>"

# Update user endpoint
curl -X POST "http://localhost:8080/api/user/update/123" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"username": "new_username"}'

```

## api field validation

> jzero integrates [https://github.com/go-playground/validator](https://github.com/go-playground/validator) by default for field validation

```shell {4}
syntax = "v1"

type CreateRequest {
    name string `json:"name" validate:"gte=2,lte=30"` // name
}
```

## Group types folder by go_package

:::important go_package option, referenced from proto files, can group generated message structs

Similarly in api files, go_package option can group defined type-generated structs

Two major advantages:
1. Avoid default generated types/types.go explosion

2. Improve development experience, type names in different groups don't conflict
:::

```shell {3,4,5,6}
syntax = "v1"

info (
	go_package: "version"
)
```

## Merge handlers of same group into one file

```shell {4}
@server (
	prefix:          /api/v1
	group:           system/user
	compact_handler: true
)
service simpleapi {
	@handler GetUserHandler
	get /system/user/getUser (GetUserRequest) returns (GetUserResponse)

	@handler DeleteUserHandler
	get /system/user/deleteUser (DeleteUserRequest) returns (DeleteUserResponse)
}
```
