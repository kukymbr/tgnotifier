package msgproc

// MessageProcessor is a tool to process messages before sending.
type MessageProcessor interface {
	// The Process receives the message text and returns it after processing.
	Process(msg string) string
}
