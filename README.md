# validator
---
<div align="center">

[![Go Report Card](https://goreportcard.com/badge/lydianpay/validator)](https://goreportcard.com/report/lydianpay/validator)
[![Maintainability](https://qlty.sh/gh/lydianpay/projects/validator/maintainability.svg)](https://qlty.sh/gh/lydianpay/projects/validator)
[![Code Coverage](https://qlty.sh/gh/lydianpay/projects/validator/coverage.svg)](https://qlty.sh/gh/lydianpay/projects/validator)
[![CodeQL](https://github.com/lydianpay/validator/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/lydianpay/validator/actions/workflows/github-code-scanning/codeql)

</div>

A small, zero-dependency struct validator for Go. Tag struct fields with
`validator:"..."` and call `ValidateStruct` to get back a single error
describing every field that failed.

## Installation & Usage

1. With [Go](https://go.dev/doc/install) installed, add the dependency:
```shell
go get -u github.com/lydianpay/validator
```
2. Tag a struct and validate it:
```go
package main

import "github.com/lydianpay/validator"

type Account struct {
	Email   string `validator:"required;email"`
	Age     int    `validator:"gte=18"`
	Species string `validator:"oneof=cat,dog,raccoon,possum"`
}

func main() {
	if err := validator.ValidateStruct(Account{Age: 17}); err != nil {
		// err lists every failing field
	}
}
```

## Supported rules

| Tag | Meaning |
|-----|---------|
| `required` | field must be non-zero |
| `email` | valid email address (string fields) |
| `exact=value` | string equals the given value |
| `oneof=a,b,c` | value is one of the listed options |
| `eq=` `gt=` `gte=` `lt=` `lte=` | numeric comparison (int, float64) |

Combine rules on a field with `;`, e.g. `validator:"required;gte=18"`.
