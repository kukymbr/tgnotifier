package buildinfo

import (
	"fmt"
	"github.com/kukymbr/tgnotifier/internal/tgnotifier"
)

// PrintVersion prints a tgnotifier version with an additional build info.
func PrintVersion() {
	fmt.Printf("tgnotifier version: %s\n\n", tgnotifier.Version)
	fmt.Printf("With gRPC server: %t\n", WithGRPCServer())
	fmt.Printf("With HTTP server: %t\n", WithHTTPServer())
}
