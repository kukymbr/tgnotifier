package msgproc

import "strings"

// NewMessageProcessor returns new MessageProcessor instance.
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{}
}

// MessageProcessor is a tool to process messages before sending.
type MessageProcessor struct{}

// Process a message.
func (mp *MessageProcessor) Process(msg string) string {
	msg = strings.TrimSpace(msg)

	return msg
}
