# REST API Patterns

## Core Architecture

### Three-Layer Pattern

jzero REST APIs follow a strict three-layer architecture:

1. **Handler Layer** (`internal/handler/`) - HTTP concerns only
2. **Logic Layer** (`internal/logic/`) - Business logic implementation
3. **Service Context** (`internal/svc/`) - Dependency injection

```
HTTP Request → Handler → Logic → External Services/Database
                  ↓
            Service Context (dependencies)
```

## Request/Response Types

### ✅ Correct Pattern

Define clear types with proper validation tags:

```go
// API definition (.api file)
type (
    CreateUserRequest {
        Name     string `json:"name" validate:"required,min=2,max=50"`
        Email    string `json:"email" validate:"required,email"`
        Age      int    `json:"age" validate:"required,gte=18,lte=120"`
        Password string `json:"password" validate:"required,min=8"`
    }

    CreateUserResponse {
        Id      int64  `json:"id"`
        Message string `json:"message"`
    }

    GetUserRequest {
        Id int64 `path:"id" validate:"required,gt=0"`
    }

    GetUserResponse {
        Id    int64  `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
        Age   int    `json:"age"`
    }

    ListUsersRequest {
        Page     int    `form:"page,default=1" validate:"gte=1"`
        PageSize int    `form:"page_size,default=10" validate:"gte=1,lte=100"`
        Keyword  string `form:"keyword,optional"`
    }

    ListUsersResponse {
        Total int64       `json:"total"`
        Users []UserInfo  `json:"users"`
    }

    UserInfo {
        Id    int64  `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }
)
```

**Tag Reference:**
- `json` - JSON field name
- `path` - Path parameter (e.g., `/users/:id`)
- `form` - Query parameter or form data
- `header` - HTTP header
- `validate` - Validation rules
- `optional` - Field is optional
- `default` - Default value

## Complete API Definition Example

```api
// user.api
info(
    title: "User API"
    desc: "User management API"
    author: "jzero"
    version: "v1"
)

type (
    CreateUserRequest {
        Name     string `json:"name" validate:"required"`
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required,min=8"`
    }

    CreateUserResponse {
        Id int64 `json:"id"`
    }

    GetUserRequest {
        Id int64 `path:"id"`
    }

    GetUserResponse {
        Id    int64  `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }

    UpdateUserRequest {
        Id   int64  `path:"id"`
        Name string `json:"name,optional"`
    }

    DeleteUserRequest {
        Id int64 `path:"id"`
    }
)

@server(
    prefix: /api/v1
    group: user
    middleware: Auth
)
service user-api {
    @doc "Create a new user"
    @handler CreateUser
    post /users (CreateUserRequest) returns (CreateUserResponse)

    @doc "Get user by ID"
    @handler GetUser
    get /users/:id (GetUserRequest) returns (GetUserResponse)

    @doc "Update user"
    @handler UpdateUser
    put /users/:id (UpdateUserRequest)

    @doc "Delete user"
    @handler DeleteUser
    delete /users/:id (DeleteUserRequest)
}
```

## When to Use This Pattern

Use the standard three-layer REST pattern for:
- CRUD APIs
- RESTful web services
- API gateways
- Backend-for-frontend (BFF) services
- Microservice APIs
