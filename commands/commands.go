package commands

import (
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"strings"
)

func ping(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error {
	utils.SendMessage(rtm, strings.Join(fields, " "), m.Channel)
	return nil
}

func init() {
	addHandler("ping", HandlerFunc(ping))
}

func HandleMessage(rtm *slack.RTM, m *slack.MessageEvent) {
	fields := strings.Fields(m.Text)
	cmd := strings.ToLower(fields[0])
	if handler, ok := handlers[cmd]; ok {
		if err := handler.HandleCommand(rtm, fields[1:], m); err != nil {
			utils.Log.Error("Failed to run command '%s': %s", cmd, err.Error())
		}
	}
}
