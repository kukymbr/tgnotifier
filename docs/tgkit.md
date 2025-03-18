# tgkit

The `tgkit` is a library that covers needs of the `tgnotifier`
with a Telegram API communications.

If this little model of the Telegram API somehow suitable to your case, 
fill free to use it as a library.
However, if you require comprehensive mappings to the Telegram API, 
it might be more suitable to use one of the fancy libraries that provide these.

## Usage

Use `go get` to add it to the project:

```shell
go get github.com/kukymbr/tgnotifier
```

Send message example:

```go
package mycoolapp

import "github.com/kukymbr/tgnotifier/pkg/tgkit"

func SendCoolMessage() error {
	// Create bot to send message from:
	bot, err := tgkit.NewBot("bot12345:bot_super_secret_token")
	if err != nil {
		return err
	}

	// Create chat ID to send message to:
	chatID := tgkit.ChatID("@chatForCoolMessages")

	// Create client:
	client := tgkit.NewDefaultClient()

	// Send the message:
	_, err = client.SendMessage(bot, tgkit.TgMessageRequest{
		ChatID: chatID,
		Text:   "😎 my unbelievable cool message",
	})

	return err
}
```

See the godoc for the reference: [pkg.go.dev](https://pkg.go.dev/github.com/kukymbr/tgnotifier@v1.0.1/pkg/tgkit)