package types

const (
	DefaultBotName  BotName  = "__default__"
	DefaultChatName ChatName = "__default__"
)

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
	return "string"
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
	return "string"
}

// UserName is a user's name from the configs.
type UserName string

func (n *UserName) String() string {
	return string(*n)
}

func (n *UserName) Set(val string) error {
	*n = UserName(val)

	return nil
}

func (n *UserName) Type() string {
	return "string"
}
