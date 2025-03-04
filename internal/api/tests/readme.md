# tgnotifier API tests

## gRPC tests

The gRPC test builds and runs the tgnotifer's gRPC server and tries to send a message through it.

To run these tests:
1. create a test config file in the [testdata/configs](testdata/configs) dir
   with a `.tgnotifier.grpc_tests.yml` filename 
   (see the [.tgnotifier.grpc_tests.example.yml](testdata/configs/.tgnotifier.grpc_tests.example.yml) example file);
2. define the next environment variables:
    * `TEST_BOT_NAME`: bot name from the `.tgnotifier.grpc_tests.yml` config
    * `TEST_CHAT_NAME`: chat name from the `.tgnotifier.grpc_tests.yml` config 
      (**message will be sent to this chat for real**)
    * `TEST_GRPC_HOST` (optional): host of the gRPC server to test (`localhost` is default)
    * `TEST_GRPC_PORT` (optional): port of the gRPC server (`50051` is default)
3. run the `go test` with a `grpc_tests` build tag:
    * `go test ./internal/api/tests/... -v -tags grpc_tests -count 1`
    * or `make test_grpc`