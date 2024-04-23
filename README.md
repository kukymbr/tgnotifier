# tgnotifier

[![License](https://img.shields.io/github/license/kukymbr/tgnotifier.svg)](https://github.com/kukymbr/tgnotifier/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/kukymbr/tgnotifier.svg)](https://github.com/kukymbr/tgnotifier/releases/latest)

The `tgnotifier` is a tool to send a notification messages
via the Telegram API.

## Installation

For now, only the `go install` is supported to install the `tgnotifier`:

```shell
go install github.com/kukymbr/tgnotifier/cmd/tgnotifier@latest
```

## CLI tool usage

```text
Supports environment variables:
- TGNOTIFIER_DEFAULT_BOT: bot identity used if no --bot flag is provided;
- TGNOTIFIER_DEFAULT_CHAT: chat ID used if no --chat flag is provided.

Usage:
  tgnotifier [flags]

Flags:
      --bot bot name           Bot name to send message from (defined in config); if not set, the bot from the TGNOTIFIER_DEFAULT_BOT env var will be used
      --chat chat name         Chat name to send message to (defined in config); if not set, the chat ID from the TGNOTIFIER_DEFAULT_CHAT env var will be used
      --config string          Path to a config file
      --disable-notification   Disable message sound notification
  -h, --help                   help for tgnotifier
      --parse-mode string      Parse mode (MarkdownV2|HTML)
      --protect-content        Protect message content from copying and forwarding
      --text string            Message text
```

### Configuration

The `tgnotifier` could have a configuration file to use multiple bots and chats.
See the [.tgnotifier.example.yml](.tgnotifier.example.yml) for an available values.

To run `tgnotifier` without the config file, define the env vars with default bot credentials and chat ID:

```shell
export TGNOTIFIER_DEFAULT_BOT="bot12345:bot-token"
export TGNOTIFIER_DEFAULT_CHAT="-12345"
```

### Command execution examples

Send a "Hello, World!" message from the default bot to the default chat:

```shell
export TGNOTIFIER_DEFAULT_BOT="bot12345:bot-token"
export TGNOTIFIER_DEFAULT_CHAT="-12345"

tgnotifier --text="Hello, World!"
```

Send a "Hello, World!" message from the `first_bot` to the `main_chat`:

```shell
tgnotifier --bot=first_bot --chat=main_chat --text="Hello, World!" 
```

Send a "Hello, World!" message from the `second_bot` to the `main_chat` with no sound notification:

```shell
tgnotifier --bot=second_bot --chat=main_chat --text="Hello, World!" --disable-notification=true
```

Send a "Hello, World!" message from the `another_bot` to the `another_chat` using the non-default config file:

```shell
tgnotifier --config="/path/to/another_config.yaml" --bot=another_bot --chat=another_chat --text="Hello, World!" 
```

## TODO

- [x] Add users' IDs to the config file to mention people in messages in format `@username`.

## License

[MIT licensed](LICENSE).