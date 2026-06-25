# validator — CLAUDE.md

## Overview
A small, zero-dependency struct validator for Go. Public module
`github.com/lydianpay/validator`, Go 1.24.5.

## API
- `ValidateStruct(input any) error` — reflects over a struct (or pointer to
  one), applies each field's `validator:"..."` tag, and returns a single error
  listing every failure, or nil.
- `DescribeBase(input any) error` — debug helper that prints a table of each
  field, its type, and its validator tag.

## Tag vocabulary (constants.go)
`required`, `email`, `exact=`, `oneof=a,b,c`, and numeric comparisons
`eq=`/`gt=`/`gte=`/`lt=`/`lte=` (int and float64). Combine rules on one field
with `;`; separate `oneof` choices with `,`. An empty/absent tag is a no-op; a
malformed tag is reported as an error.

## Layout
- `constants.go` — tag names and separators.
- `funcs.go` — `ValidateStruct` and the per-rule `validateField` dispatch.
- `util.go` — reflection/tag-parsing helpers, number coercion, `DescribeBase`.

## Testing
Native `testing` package, no external dependencies.

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```
