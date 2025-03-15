# tgnotifier GUI

The `tgnotifierui` app is a GUI wrapper for a `tgnotifier`.

<img src="../../docs/gui.png" width="600">

## Installation

To install the `tgnotifierui`, download the binary for your OS:

* Ubuntu/Debian: [tgnotifierui_v1.0.0_ubuntu-latest](https://github.com/kukymbr/tgnotifier/releases/download/v1.0.0/tgnotifierui_v1.0.0_ubuntu-latest)
* macOS: [tgnotifierui_v1.0.0_macOS-latest](https://github.com/kukymbr/tgnotifier/releases/download/v1.0.0/tgnotifierui_v1.0.0_macOS-latest)
* Windows: [tgnotifierui_v1.0.0_windows-latest.exe](https://github.com/kukymbr/tgnotifier/releases/download/v1.0.0/tgnotifierui_v1.0.0_windows-latest.exe)

## Starting the GUI

To run the `tgnotifierui`:

1. Follow the steps from the [Configuration](../../README.md#configuration) 
   section in the main doc to create a configuration file.
2. Run the app:
   * from the console if you need to specify a path to a config file:
     ```shell
      # If you need to specify a path to a config file:
      tgnotifierui --config=/path/to/.tgnotifier.yml
     ```
   * or just double-click it if config file is in user's home or config dir
     or if env vars are used.

## Usage

Bots and chats from a configuration file are loaded into the UI's dropdowns.

Select a bot in the `Bot` dropdown and send one of the bot-related requests by clicking the buttons:

* `Check bot...` to send a `getMe` Telegram request;
* `Get updates...` to send a `getUpdates` Telegram request.

Select a chat to send message to, fill the `Message` field
and send a message by clicking to the `Send it now!` button.