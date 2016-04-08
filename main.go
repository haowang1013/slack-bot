package main

import (
	"fmt"
	"github.com/haowang1013/slack-bot/commands"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"os"
)

const (
	channelID    = "C0E0ST6L9"
	tokenEnvName = "SLACK_API_TOKEN"
)

func main() {
	token := os.Getenv(tokenEnvName)
	if len(token) == 0 {
		panic(fmt.Sprintf("Failed to get slack api token from environment variable '%s'", tokenEnvName))
	}

	api := slack.New(token)

	utils.ListGroups(api)
	utils.ListChannels(api)

	//api.SetDebug(true)
	rtm := api.NewRTM()
	commands.SetRTM(rtm)
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				utils.SendMessage(rtm, "I'm connected", channelID)
			case *slack.MessageEvent:
				commands.HandleMessage(ev)
			case *slack.LatencyReport:
				utils.Log.Debugf("Current latency: %v", ev.Value)
			case *slack.InvalidAuthEvent:
				utils.Log.Error("Invalid credentials")
				break
			}
		}
	}
}
