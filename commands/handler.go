package commands

import (
	"errors"
	"fmt"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"strings"
)

type CommandHandler interface {
	HandleCommand(fields []string, m *slack.MessageEvent) error
}

type HandlerFunc func(fields []string, m *slack.MessageEvent) error

func (f HandlerFunc) HandleCommand(fields []string, m *slack.MessageEvent) error {
	return f(fields, m)
}

type handleCollection map[string]CommandHandler

func (h *handleCollection) HandleCommand(fields []string, m *slack.MessageEvent) error {
	if len(fields) == 0 {
		return errors.New("Missing param")
	}

	cmd := fields[0]
	if handler, ok := (*h)[cmd]; ok {
		return handler.HandleCommand(fields[1:], m)
	} else {
		return errors.New(fmt.Sprintf("Unrecognized param: %s", cmd))
	}
}

func (h *handleCollection) addHandler(name string, handler CommandHandler) error {
	name = strings.ToLower(name)
	if h, ok := (*h)[name]; ok {
		msg := fmt.Sprintf("Handler for command '%s' is already registered: %v", name, h)
		utils.Log.Error(msg)
		return errors.New(msg)
	}
	(*h)[name] = handler
	return nil
}

var rootHandler = handleCollection{}
