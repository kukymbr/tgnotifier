package config

import (
	"github.com/kukymbr/tgnotifier/internal/types"
)

type configDTO struct {
	Bots  map[types.BotName]string  `json:"bots" yaml:"bots"`
	Chats map[types.ChatName]string `json:"chats" yaml:"chats"`

	DefaultBot  types.BotName  `json:"default_bot" yaml:"default_bot"`
	DefaultChat types.ChatName `json:"default_chat" yaml:"default_chat"`

	SilenceSchedule []silenceScheduleItem `json:"silence_schedule" yaml:"silence_schedule"`

	Replaces map[string]string `json:"replaces" yaml:"replaces"`

	GRPC serverConfigDTO `json:"grpc" yaml:"grpc"`
	HTTP serverConfigDTO `json:"http" yaml:"http"`
}

type silenceScheduleItem struct {
	From string `json:"from" yaml:"from"`
	To   string `json:"to" yaml:"to"`
}

type serverConfigDTO struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}
