<img align="right" width="125" src="assets/tgnotifier.png" alt="image with a gopher on a telegram plane">

# tgnotifier

[![License](https://img.shields.io/github/license/kukymbr/tgnotifier.svg)](https://github.com/kukymbr/tgnotifier/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/kukymbr/tgnotifier.svg)](https://github.com/kukymbr/tgnotifier/releases/latest)
[![GoDoc](https://godoc.org/github.com/kukymbr/tgnotifier?status.svg)](https://godoc.org/github.com/kukymbr/tgnotifier)
[![GoReport](https://goreportcard.com/badge/github.com/kukymbr/tgnotifier)](https://goreportcard.com/report/github.com/kukymbr/tgnotifier)

The `tgnotifier` is a tool to send a notification messages
via the Telegram API.

## Okay, but why?

The main feature of the `tgnotifier` is in the configuration, 
you could define bots and chats and just send the messages, just like this: 

```shell
tgnotifier send --text="😎 my cool message"
```

or like this:

```shell
tgnotifier send --text="🦖 my the most coolest message" --bot="cool_messages_bot" --chat="chat_for_cool_messages"
```

One of the cases when you may want to use `tgnotifier` is an alerting from the CI/CD or some other services
via the single entrypoint without exposing the bot credentials in every service or configuration. 
The gRPC server is also presented in the `tgnotifier` to have it running as a service.

## Installation

To install the latest release of the `tgnotifier`, 
download the archive with the binary for your OS and unpack it somewhere inside the PATH.

* Ubuntu/Debian: [tgnotifier_v0.7.0_ubuntu-latest.zip](https://github.com/kukymbr/tgnotifier/releases/download/v0.7.0/tgnotifier_v0.7.0_ubuntu-latest.zip)
* Windows: [tgnotifier_v0.7.0_windows-latest.zip](https://github.com/kukymbr/tgnotifier/releases/download/v0.7.0/tgnotifier_v0.7.0_windows-latest.zip)

Installation on Ubuntu example:

```shell
wget https://github.com/kukymbr/tgnotifier/releases/download/v0.7.0/tgnotifier_v0.7.0_ubuntu-latest.zip
unzip tgnotifier_v0.7.0_ubuntu-latest.zip -d /usr/local/bin/
tgnotifier --version
```

> The `--version` flag is supported since the v0.3.2 version.

### Compile from sources

To install `tgnotifier` from the source, use the `go install` command:

```shell
go install github.com/kukymbr/tgnotifier/cmd/tgnotifier@v0.7.0
```

<details>
  <summary><b>Custom build without some components</b></summary>

There is a possibility to compile a custom build of the `tgnotifier`.
Golang 1.23 and above is required.

Available since v0.7.0.

```shell
# Clone the repository:
git clone https://github.com/kukymbr/tgnotifier.git && cd tgnotifier

# Checkout tag you want to compile:
git checkout v0.7.0

# To build the tgnotifier with all components:
make build

# To build the tgnotifier without the gRPC server:
make build_without_gprc

# To build the tgnotifier without the HTTP server:
make build_without_http

# To build the tgnotifier without the gRPC and HTTP server both:
make build_without_servers
```
</details>

You could check type of the build by running the `tgnotifier version` command:

```shell
./tgnotifier version

# tgnotifier version: v0.7.0
#
# With gRPC server: true
# With HTTP server: true
```

### Docker

There is a Docker configuration for a `tgnotifier`.

Pure `docker run` call example: 

```shell
docker run --env TGNOTIFIER_DEFAULT_BOT=bot12345:bot_token --env TGNOTIFIER_DEFAULT_CHAT=-12345 --rm ghcr.io/kukymbr/tgnotifier:0.5.0 send --text="what's up?"
```

See the [docker/compose.yml](docker/compose.yml) file for a docker compose usage example.

## Configuration

The `tgnotifier` could have a configuration file to use multiple bots and chats.
See the [.tgnotifier.example.yml](.tgnotifier.example.yml) for an available values.

Defining the bots, who can send messages via the `tgnotifier`:

```yaml
bots:
  first_bot: "12345:FIRST_BOT_TOKEN"
  second_bot: "bot54321:SECOND_BOT_TOKEN"
```

Defining the chat IDs, where `tgnotifier` can send messages to:

```yaml
chats:
  main_chat: "-12345"
  secondary_chat: "@my_test_channel"
```

To use a program without a `bot` or a `chat` argument,
define a `default_bot` and  `default_chat` values in the config file:

```yaml
default_bot: "first_bot"
default_chat: "main_chat"
```

Defining the schedule when messages are sent without a sound notification:

```yaml
silence_schedule:
  - from: 22:00
    to: 9:00
```

Defining substrings to replace in the message texts:

```yaml
replaces:
  "@test": "[Test](tg://user?id=123456789)"
  "FYI": "for your information"
```

To run `tgnotifier` without the config file at all,
define the env vars with default bot credentials and chat ID:

```shell
export TGNOTIFIER_DEFAULT_BOT="bot12345:bot-token"
export TGNOTIFIER_DEFAULT_CHAT="-12345"
```

## CLI tool usage

```text
Usage:
  tgnotifier send [flags]

Flags:
      --bot string                      Bot name to send message from (defined in config); if not set, the default_bot directive or the bot from the TGNOTIFIER_DEFAULT_BOT env var will be used
      --chat string                     Chat name to send message to (defined in config); if not set, the default_chat directive or the chat ID from the TGNOTIFIER_DEFAULT_CHAT env var will be used
      --config string                   Path to a config file
      --debug                           Enable the debug mode
      --disable-notification            Disable message sound notification
  -h, --help                            help for send
      --parse-mode message parse mode   Parse mode (MarkdownV2|HTML)
      --protect-content                 Protect message content from copying and forwarding
      --text string                     Message text
  -v, --version                         version for send
```

### Command execution examples

Send a "Hello, World!" message from the default bot to the default chat:

```shell
export TGNOTIFIER_DEFAULT_BOT="bot12345:bot-token"
export TGNOTIFIER_DEFAULT_CHAT="-12345"

tgnotifier --text="Hello, World!"
```

Send a "Hello, World!" message from the `first_bot` to the `main_chat`, mentioning the `JohnDoe` user:

```shell
tgnotifier --bot=first_bot --chat=main_chat --text="Hello, World and @JohnDoe!"
```

Send a "Hello, World!" message from the `second_bot` to the `main_chat` with no sound notification:

```shell
tgnotifier --bot=second_bot --chat=main_chat --text="Hello, World!" --disable-notification=true
```

Send a "Hello, World!" message from the `another_bot` to the `another_chat` using the non-default config file:

```shell
tgnotifier --config="/path/to/another_config.yaml" --bot=another_bot --chat=another_chat --text="Hello, World!" 
```

## gRPC server

The tgnotifier could be started as an gRPC server. To run the server, use the `tgnotifier grpc` command:

```text
Usage:
  tgnotifier grpc [flags]

Flags:
      --config string   Path to a config file
      --debug           Enable the debug mode
  -h, --help            help for grpc
  -v, --version         version for grpc
```

See the [tgnotifier.proto](api/grpc/tgnotifier.proto) for an API contract.

## HTTP server

The tgnotifier could be started as an HTTP server too. 
To run the server, use the `tgnotifier http` command:

```text
Usage:
  tgnotifier http [flags]

Flags:
      --config string   Path to a config file
      --debug           Enable the debug mode
  -h, --help            help for grpc
  -v, --version         version for grpc
```

See the [openapi.yaml](api/http/openapi.yaml) for an API contract.

## `tgkit` library

Yes, there also a library to communicate with a Telegram API.
See its doc for more info: [pkg/tgkit/readme.md](pkg/tgkit/readme.md).

## TODO

- [ ] `tgkit`: support all the Telegram response fields. 
- [ ] Predefined messages with templates and i18n. 
- [x] HTTP server.
- [x] Optional user's config file in the home dir.
- [x] gRPC server.
- [x] Docker configuration.
- [x] Replace map in config.
- [x] Define default bot & chat in config file.
- [x] Silence schedule in config.

## License

[MIT licensed](LICENSE).