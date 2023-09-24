//go:build tools
// +build tools

package tools

// TODO: in go 1.22, replace this with go tools

import (
	_ "github.com/atombender/go-jsonschema/cmd/gojsonschema"
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
