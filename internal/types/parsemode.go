package types

const (
	ParseModeDefault   ParseMode = ""
	ParseModeMarkdown2 ParseMode = "MarkdownV2"
	ParseModeHTML      ParseMode = "HTML"
)

// ParseMode is a message parse mode.
type ParseMode string

func (m *ParseMode) String() string {
	return string(*m)
}

func (m *ParseMode) Set(val string) error {
	*m = ParseMode(val)

	return nil
}

func (m *ParseMode) Type() string {
	return "message parse mode"
}
