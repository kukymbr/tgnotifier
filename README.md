# tgnotifier

[![License](https://img.shields.io/github/license/kukymbr/tgnotifier.svg)](https://github.com/kukymbr/tgnotifier/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/kukymbr/tgnotifier.svg)](https://github.com/kukymbr/tgnotifier/releases/latest)

The `tgnotifier` is a tool to send a notification messages
via the Telegram API.

## Installation

For now, only the `go install` is supported to install the `tgnotifier`:

```shell
go install github.com/kukymbr/tgnotifier@latest
```

## CLI tool usage

To run the `tgnotifier`, the configuration file is required.
See the [.tgnotifier.example.yml](.tgnotifier.example.yml) for an available values.

1. Create the configuration file, by default it is `.tgnotifier.yml`.
2. Produce the bot credentials and the chat IDs into the config file.
3. Run the `tgnotifier`.

```text
Usage:
  tgnotifier [flags]

Flags:
      --bot bot name           Bot name to send message from (defined in config)
      --chat chat name         Chat name to send message to (defined in config)
      --config string          Path to a config file (default ".tgnotifier.yml")
      --disable-notification   Disable message sound notification
  -h, --help                   help for tgnotifier
      --parse-mode string      Parse mode (MarkdownV2|HTML)
      --protect-content        Protect message content from copying and forwarding
      --text string            Message text
```

### Command execution examples

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

## License

[MIT licensed](LICENSE).