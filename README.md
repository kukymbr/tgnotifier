<img align="right" width="125" src="assets/tgnotifier.png" alt="image with a gopher on a telegram plane">

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
      --bot string                      Bot name to send message from (defined in config); if not set, the bot from the TGNOTIFIER_DEFAULT_BOT env var will be used
      --chat string                     Chat name to send message to (defined in config); if not set, the chat ID from the TGNOTIFIER_DEFAULT_CHAT env var will be used
      --config string                   Path to a config file
      --debug                           Enable the debug mode
      --disable-notification            Disable message sound notification
  -h, --help                            help for tgnotifier
      --parse-mode message parse mode   Parse mode (MarkdownV2|HTML)
      --protect-content                 Protect message content from copying and forwarding
      --text string                     Message text
```

### Configuration

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

To run `tgnotifier` without the config file at all, 
define the env vars with default bot credentials and chat ID:

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

## TODO

- [x] Define default bot & chat in config file.
- [ ] Docker configuration.

## License

[MIT licensed](LICENSE).