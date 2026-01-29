# Proto Field Validation

## Overview

jzero supports proto field validation using [protovalidate](https://buf.build/docs/protovalidate/). This allows you to define validation rules directly in your proto files, ensuring data integrity at the RPC layer.

## Validation with protovalidate

jzero uses the `buf/validate/validate.proto` package to provide comprehensive field validation capabilities.

## Basic Example

```protobuf
syntax = "proto3";

package versionpb;

import "buf/validate/validate.proto";

option go_package = "./pb/versionpb";

message VersionRequest {
  int32 id = 1 [
    (buf.validate.field).cel = {
      id: "id.length"
      message: "id must be greater than 0 and less than 100000"
      expression: "this > 0 && this < 100000"
    }
  ];
}
```

## Common Validation Patterns

### String Validation

```protobuf
message CreateUserRequest {
  string username = 1 [
    (buf.validate.field).string.min_len = 3,
    (buf.validate.field).string.max_len = 32,
    (buf.validate.field).string.pattern = "^[a-zA-Z0-9_]+$"
  ];

  string email = 2 [
    (buf.validate.field).string.email = true
  ];

  string phone = 3 [
    (buf.validate.field).string.min_len = 10,
    (buf.validate.field).string.max_len = 15
  ];
}
```

### Numeric Validation

```protobuf
message ProductRequest {
  int32 quantity = 1 [
    (buf.validate.field).int32.greater_than = 0,
    (buf.validate.field).int32.less_than = 10000
  ];

  double price = 2 [
    (buf.validate.field).double.greater_than = 0.0,
    (buf.validate.field).double.less_than = 1000000.0
  ];

  uint64 id = 3 [
    (buf.validate.field).uint64.const = 1
  ];
}
```

### CEL (Common Expression Language) Validation

For complex validation logic, use CEL expressions:

```protobuf
message OrderRequest {
  int32 item_id = 1 [
    (buf.validate.field).cel = {
      id: "item_id.range"
      message: "item_id must be between 1 and 99999"
      expression: "this > 0 && this < 100000"
    }
  ];

  string status = 2 [
    (buf.validate.field).cel = {
      id: "status.valid"
      message: "status must be pending, processing, or completed"
      expression: "this in ['pending', 'processing', 'completed']"
    }
  ];

  int64 start_time = 3 [
    (buf.validate.field).cel = {
      id: "start_time.future"
      message: "start_time must be in the future"
      expression: "this > timestamp.now()"
    }
  ];
}
```

### Required Fields

```protobuf
message RequiredFieldsExample {
  string name = 1 [
    (buf.validate.field).required = true
  ];

  int32 age = 2 [
    (buf.validate.field).required = true
  ];
}
```

### Message Validation

Validate nested message structures:

```protobuf
message Address {
  string street = 1 [(buf.validate.field).required = true];
  string city = 2 [(buf.validate.field).required = true];
  string zip_code = 3 [
    (buf.validate.field).string.pattern = "^\\d{5}$"
  ];
}

message CreateUserRequest {
  string name = 1 [(buf.validate.field).required = true];
  Address address = 2 [
    (buf.validate.field).required = true,
    (buf.validate.field).message.required = true
  ];
}
```

### Repeated Field Validation

```protobuf
message BatchRequest {
  repeated string tags = 1 [
    (buf.validate.field).repeated.min_items = 1,
    (buf.validate.field).repeated.max_items = 10,
    (buf.validate.field).repeated.unique = true
  ];

  repeated int64 ids = 2 [
    (buf.validate.field).repeated.min_items = 1,
    (buf.validate.field).repeated.max_items = 100
  ];
}
```

### Map Validation

```protobuf
message MetadataRequest {
  map<string, string> metadata = 1 [
    (buf.validate.field).map.min_pairs = 1,
    (buf.validate.field).map.max_pairs = 10
  ];
}
```

### Enum Validation

```protobuf
enum Status {
  STATUS_UNKNOWN = 0;
  STATUS_ACTIVE = 1;
  STATUS_INACTIVE = 2;
  STATUS_PENDING = 3;
}

message UpdateStatusRequest {
  Status status = 1 [
    (buf.validate.field).enum = {
      defined_only = true,
      not_in: [0]  // Exclude UNKNOWN
    }
  ];
}
```

## Advanced Validation Patterns

### Cross-Field Validation

```protobuf
message DateRangeRequest {
  int64 start_date = 1;
  int64 end_date = 2;

  // Validate at message level
  option (buf.validate.message).cel = {
    id: "date_range.valid"
    message: "end_date must be after start_date"
    expression: "this.end_date > this.start_date"
  };
}
```

### Conditional Validation

```protobuf
message ConditionalRequest {
  string user_type = 1;
  string company_name = 2;
  string personal_name = 3;

  option (buf.validate.message).cel = {
    id: "conditional.names"
    message: "company_name required for business, personal_name for individual"
    expression: "this.user_type == 'business' ? this.company_name != '' : this.personal_name != ''"
  };
}
```

## Built-in Constraints

### String Constraints

- `min_len`: Minimum length
- `max_len`: Maximum length
- `pattern`: Regex pattern
- `prefix`: Must start with
- `suffix`: Must end with
- `contains`: Must contain
- `email`: Must be valid email format
- `hostname`: Must be valid hostname
- `ip`: Must be valid IP address
- `uuid`: Must be valid UUID format
- `uri`: Must be valid URI

### Numeric Constraints

- `greater_than`: Must be > value
- `less_than`: Must be < value
- `greater_than_or_equal`: Must be >= value
- `less_than_or_equal`: Must be <= value
- `const`: Must equal value
- `in`: Must be in set of values
- `not_in`: Must not be in set of values

### Repeated Constraints

- `min_items`: Minimum number of items
- `max_items`: Maximum number of items
- `unique`: All items must be unique

## Error Handling

When validation fails, jzero returns detailed error messages:

```go
type ValidationError struct {
    Field   string
    Message string
    Constraint string
}
```

Example error response:
```json
{
  "code": 400,
  "message": "validation failed",
  "details": [
    {
      "field": "username",
      "message": "username length must be at least 3 characters",
      "constraint": "min_len"
    }
  ]
}
```

## Best Practices

### ✅ Do's

```protobuf
// Use descriptive error messages
int32 age = 1 [
  (buf.validate.field).cel = {
    id: "age.range"
    message: "age must be between 18 and 120"
    expression: "this >= 18 && this <= 120"
  }
];

// Use CEL for complex logic
string password = 1 [
  (buf.validate.field).cel = {
    id: "password.strength"
    message: "password must contain uppercase, lowercase, and number"
    expression: "this.matches('^[A-Z]') && this.matches('^[a-z]') && this.matches('^[0-9]')"
  }
];

// Validate at both field and message level when needed
message DateRange {
  int64 start = 1 [(buf.validate.field).required = true];
  int64 end = 2 [(buf.validate.field).required = true];

  option (buf.validate.message).cel = {
    id: "range.valid"
    message: "end must be after start"
    expression: "this.end > this.start"
  };
}
```

### ❌ Don'ts

```protobuf
// DON'T: Skip validation for user input
string email = 1;  // ❌ No validation

// DON'T: Use vague error messages
int32 age = 1 [
  (buf.validate.field).cel = {
    id: "age.check"
    message: "invalid"  // ❌ Not helpful
    expression: "this > 0"
  }
];

// DON'T: Over-complicate simple validation
string name = 1 [
  (buf.validate.field).cel = {
    id: "name.complex"
    message: "name validation"
    expression: "this.size() >= 1 && this.size() <= 100"  // ❌ Use min_len/max_len instead
  }
];
```

## Complete Example

```protobuf
syntax = "proto3";

package userpb;

import "buf/validate/validate.proto";
import "google/api/annotations.proto";

option go_package = "./pb/userpb";

message CreateUserRequest {
  string username = 1 [
    (buf.validate.field).required = true,
    (buf.validate.field).string.min_len = 3,
    (buf.validate.field).string.max_len = 32,
    (buf.validate.field).string.pattern = "^[a-zA-Z0-9_]+$"
  ];

  string email = 2 [
    (buf.validate.field).required = true,
    (buf.validate.field).string.email = true
  ];

  int32 age = 3 [
    (buf.validate.field).required = true,
    (buf.validate.field).int32.greater_than = 0,
    (buf.validate.field).int32.less_than = 150
  ];

  repeated string tags = 4 [
    (buf.validate.field).repeated.max_items = 10,
    (buf.validate.field).repeated.unique = true
  ];
}

message CreateUserResponse {
  int32 id = 1;
  string username = 2;
}

service UserService {
  rpc CreateUser(CreateUserRequest) returns(CreateUserResponse) {
    option (google.api.http) = {
      post: "/api/v1/user/create"
      body: "*"
    };
  };
}
```

## Resources

- [protovalidate Documentation](https://buf.build/docs/protovalidate/)
- [CEL Language Specification](https://github.com/google/cel-spec)
- [jzero Documentation](https://docs.jzero.io)
