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
you could define bots and chats once and just send the messages, just like this: 

```shell
tgnotifier send --text="😎 my cool message"
```

or like this:

```shell
tgnotifier send --text="🦖 my the most coolest message" --bot="cool_messages_bot" --chat="chat_for_cool_messages"
```

One of the cases when you may want to use `tgnotifier` is an alerting from the CI/CD or some other services
via the single entrypoint without exposing the bot credentials in every service or configuration. 
The gRPC and HTTP servers are also presented in the `tgnotifier` to have it running as a service.

## Installation

To install the latest release of the `tgnotifier`, 
download the binary for your OS and put it somewhere inside the PATH.

* Ubuntu/Debian: [tgnotifier_v1.0.1_ubuntu-latest](https://github.com/kukymbr/tgnotifier/releases/download/v1.0.1/tgnotifier_v1.0.1_ubuntu-latest)
* macOS: [tgnotifier_v1.0.1_macOS-latest](https://github.com/kukymbr/tgnotifier/releases/download/v1.0.1/tgnotifier_v1.0.1_macOS-latest)
* Windows: [tgnotifier_v1.0.1_windows-latest.exe](https://github.com/kukymbr/tgnotifier/releases/download/v1.0.1/tgnotifier_v1.0.1_windows-latest.exe)

Installation on Ubuntu/macOS example:

```shell
wget https://github.com/kukymbr/tgnotifier/releases/download/v1.0.1/tgnotifier_v1.0.1_ubuntu-latest
chmod a+x tgnotifier_v1.0.1_ubuntu-latest 
mv tgnotifier_v1.0.1_ubuntu-latest /usr/local/bin/tgnotifier
tgnotifier --version
```

<details>
  <summary><b>Compile from sources</b></summary>

To install `tgnotifier` from the source, use the `go install` command:

```shell
go install github.com/kukymbr/tgnotifier/cmd/tgnotifier@v1.0.1
```
</details>

<details>
  <summary><b>Custom build without some components</b></summary>

There is a possibility to compile a custom build of the `tgnotifier`.
Golang 1.24 and above is required.

Available since v0.7.1.

```shell
# Clone the repository:
git clone https://github.com/kukymbr/tgnotifier.git && cd tgnotifier

# Checkout tag you want to compile:
git checkout v1.0.1

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

<details>
  <summary><b>Docker configuration</b></summary>

Docker image: 
[ghcr.io/kukymbr/tgnotifier:1.0.1](https://github.com/kukymbr/tgnotifier/pkgs/container/tgnotifier/375583862?tag=1.0.1).

#### Docker usage examples

Sending message without a configuration file:

```shell
docker run --env TGNOTIFIER_DEFAULT_BOT=bot12345:bot_token --env TGNOTIFIER_DEFAULT_CHAT=-12345 --rm ghcr.io/kukymbr/tgnotifier:1.0.1 send --text="what's up?"
```

Sending message with a configuration file:

```shell
docker run --volume /path/to/.tgnotifier.yml:/tgnotifier/.tgnotifier.yml --rm ghcr.io/kukymbr/tgnotifier:1.0.1 send --text="what's up?"
```

Starting a gRPC server:

```shell
docker run --volume /path/to/.tgnotifier.yml:/tgnotifier/.tgnotifier.yml --rm ghcr.io/kukymbr/tgnotifier:1.0.1 grpc
```

Using as a docker compose service:

```yaml
services:
  tgnotifier:
    image: ghcr.io/kukymbr/tgnotifier:1.0.1
    volumes:
      - /path/to/.tgnotifier.yml:/tgnotifier/.tgnotifier.yml
    ports:
      - "127.0.0.1:8080:8080"
    command: ["grpc"]
```

See the [docker](docker) directory for a full docker compose usage example.
</details>

## Configuration

The `tgnotifier` could have a configuration file to use multiple bots and chats.
See the [.tgnotifier.example.yml](.tgnotifier.example.yml) for an available values.

Use one of these ways to define a `tgnotifier` configuration:

* Create a file named `.tgnotifier.yml` in the user's home or config dir,
  or near the `tgnotifier` executable file,
  `tgnotifier` will use it if no explicit config file passed as argument:
  ```shell
  nano "$HOME/.config/.tgnotifier.yml" # Define a config values...
  tgnotifer --text="🤠" 
  ```
* Or create an YAML of JSON file in any location you want and give its path to the `tgnotifier`:
  ```shell
  tgnotifier --config="/path/to/config.yml" --text="🎉"
  ```
* Or create a configuration file and set its path to the `TGNOTIFIER_CONFIG_PATH` environment variable:
   ```shell
  export TGNOTIFIER_CONFIG_PATH=/path/to/config.yml 
  tgnotifier --text="⭐"
  ```
* Or use a minimal configuration mode (just single bot and chat) using the environment variables:
  ```shell
  export TGNOTIFIER_DEFAULT_BOT="bot12345:bot-token"
  export TGNOTIFIER_DEFAULT_CHAT="-12345"
  
  tgnotifer --text="🔥" 
  ```

### Configuration values

Put variables below to your configuration file to make `tgnnotifier` run in your way.
Examples are provided in YAML format, but JSON files are supported too.

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
  "{{ failed }}": "Pipeline job has failed" 
```

<details>
<summary><b>Some other configuration options</b></summary>

```yaml
# Telegram API client configuration.
client:
  # Telegram API requests timeout.
  timeout: 30s

# Failed requests (500 responses, timeouts, protocol errors) retrier options.
retrier:
  # Type of the retrier.
  # Available:
  #  - noop: just single attempt, no retries;
  #  - linear: fixed retry delay, `attempts` and `delay` fields are used;
  #  - progressive: increasing retry delay, `attempts`, `delay` and `multiplier` fields are used.
  type: progressive
  # Maximum count of attempts.
  attempts: 3
  # Delay between an attempts, or initial delay in case of a progressive retrier.
  delay: 500ms
  # Delay multiplier for a progressive retrier.
  multiplier: 2
```
</details>

To run `tgnotifier` without a config file at all,
define the env vars with default bot credentials and chat ID:

```shell
export TGNOTIFIER_DEFAULT_BOT="bot12345:bot-token"
export TGNOTIFIER_DEFAULT_CHAT="-12345"
```

<details>
  <summary><b>Where can I get a bot identifier and chat ID?</b></summary>

1. See the [Telegram Bot API tutorial](https://core.telegram.org/bots/tutorial#obtain-your-bot-token) 
to find out how to obtain a bot token.
2. Getting the chat ID is little more tricky, you need to add bot to a required chat, send there a message
   and request the `getUpdates` Telegram API endpoint:
   * use `curl` or browser to request this data, 
     see [this doc](https://gist.github.com/nafiesl/4ad622f344cd1dc3bb1ecbe468ff9f8a#get-chat-id-for-a-private-chat)
     for an example;
   * or use `tgnotifierui` GUI, see its doc for more info: [docs/gui.md](docs/gui.md#getting-the-chat-id). 
</details>

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

Send a "Hello, World!" message from the default bot to the default chat, defined by the environment variables:

```shell
export TGNOTIFIER_DEFAULT_BOT="bot12345:bot-token"
export TGNOTIFIER_DEFAULT_CHAT="-12345"

tgnotifier --text="Hello, World!"
```

Send a "Hello, World!" message from the default bot to the default chat, defined in the config file:

```shell
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

## GUI

If you want to test or visualize `tgnotifier` possibilities or to send couple messages via the UI, 
run the `tgnotifierui` app.

See the [GUI readme](docs/gui) for more info.

## `tgkit` library

Yes, there also a library to communicate with a Telegram API.
See its doc for more info: [docs/tgkit.md](docs/tgkit.md).

## TODO

- [ ] Predefined messages with templates and i18n. 
- [ ] `tgkit`: generate models from Telegram's swagger file.
- [x] HTTP server.
- [x] Optional user's config file in the home dir.
- [x] gRPC server.
- [x] Docker configuration.
- [x] Replace map in config.
- [x] Define default bot & chat in config file.
- [x] Silence schedule in config.

## License

[MIT licensed](LICENSE).