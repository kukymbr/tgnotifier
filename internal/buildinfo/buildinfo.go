package buildinfo

// WithGRPCServer checks if tgnotifier is built with a gRPC server.
func WithGRPCServer() bool {
	return gRPCEnabled
}

// WithHTTPServer checks if tgnotifier is built with a gRPC server.
func WithHTTPServer() bool {
	return httpEnabled
}
