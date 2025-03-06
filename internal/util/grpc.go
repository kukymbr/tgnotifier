package util

import (
	"google.golang.org/grpc/status"
)

// ErrorResponse is an alias to the status.Errorf
// to fix a `go vet` "non-constant format string" error.
// See https://github.com/grpc/grpc-go/issues/90.
var ErrorResponse = status.Errorf
