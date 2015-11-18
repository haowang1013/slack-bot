package commands

import (
	"fmt"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"strings"
	"time"
)

func pingHandler(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error {
	utils.SendMessage(rtm, strings.Join(fields, " "), m.Channel)
	return nil
}

func init() {
	addHandler("ping", HandlerFunc(pingHandler))
}

func HandleMessage(rtm *slack.RTM, m *slack.MessageEvent) {
	fields := strings.Fields(m.Text)
	cmd := strings.ToLower(fields[0])
	if handler, ok := handlers[cmd]; ok {
		begin := time.Now()
		err := handler.HandleCommand(rtm, fields[1:], m)
		if err == nil {
			duration := time.Now().Sub(begin)
			utils.Log.Info("Command '%s' executed in %.2f sec", cmd, duration.Seconds())
		} else {
			msg := fmt.Sprintf("Command '%s' failed: %s", cmd, err.Error())
			utils.Log.Error(msg)
			utils.SendMessage(rtm, msg, m.Channel)
		}
	} else {
		// command not recognized
	}
}
