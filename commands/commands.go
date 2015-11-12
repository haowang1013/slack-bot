package commands

import (
	//"encoding/json"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/kr/pretty"
	"github.com/nlopes/slack"
)

func HandleMessage(m *slack.MessageEvent) {
	if len(m.BotID) > 0 {
		// ignore bot message
		return
	}

	utils.Log.Debug("Got message: %v", pretty.Formatter(m))
}
