# HTTP spec

HTTP API specification is generated from the [tgnotifier.proto](../grpc/tgnotifier.proto) file 
using the [protoc-gen-oas](https://github.com/ogen-go/protoc-gen-oas) generator. Do not edit it by hands.

To generate an openapi file, run 

```shell
make apis
```

This will generate gRPC server stubs, openapi.yaml and HTTP server stubs.