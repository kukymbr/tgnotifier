# tgnotifier GUI

The `tgnotifierui` app is a GUI wrapper for a `tgnotifier`.

<img src="../../docs/gui.png" width="600">

## Starting the GUI

To run the `tgnotifierui`:

1. Follow the steps from the [Configuration](../../README.md#configuration) 
   section in the main doc to create a configuration file.
2. Run the app:
   ```shell
   tgnotifierui --config=/path/to/.tgnotifier.yml
   ```

## Usage

Bots and chats from a configuration file are loaded into the UI's dropdowns.

Select a bot in the `Bot` dropdown and send one of the bot-related requests by clicking the buttons:

* `Check bot...` to send a `getMe` Telegram request;
* `Get updates...` to send a `getUpdates` Telegram request.

Select a chat to send message to, fill the `Message` field
and send a message by clicking to the `Send it now!` button.