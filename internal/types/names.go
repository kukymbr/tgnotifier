package types

// BotName is a given bot name from the configs.
type BotName string

func (n *BotName) String() string {
	return string(*n)
}

func (n *BotName) Set(val string) error {
	*n = BotName(val)

	return nil
}

func (n *BotName) Type() string {
	return "bot name"
}

// ChatName is a given chat name from the configs.
type ChatName string

func (n *ChatName) String() string {
	return string(*n)
}

func (n *ChatName) Set(val string) error {
	*n = ChatName(val)

	return nil
}

func (n *ChatName) Type() string {
	return "chat name"
}