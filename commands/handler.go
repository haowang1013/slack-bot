package commands

import (
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"strings"
)

type CommandHandler interface {
	HandleCommand(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error
}

var handlers = map[string]CommandHandler{}

type HandlerFunc func(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error

func (f HandlerFunc) HandleCommand(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error {
	return f(rtm, fields, m)
}

func addHandler(name string, handler CommandHandler) {
	name = strings.ToLower(name)
	if h, ok := handlers[name]; ok {
		utils.Log.Error("Handler for command '%s' is already registered: %v", name, h)
		return
	}
	handlers[name] = handler
}
