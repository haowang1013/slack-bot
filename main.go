package main

import (
	"fmt"
	"github.com/haowang1013/slack-bot/commands"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
)

const (
	token     = "xoxb-14041206193-ruJiHzQMBkGoMwO2mThAChiw"
	channelID = "C0E0ST6L9"
)

func main() {
	api := slack.New(token)

	utils.ListGroups(api)
	utils.ListChannels(api)

	//api.SetDebug(true)
	rtm := api.NewRTM()
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
				fmt.Printf("Current latency: %v\n", ev.Value)
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break
			}
		}
	}
}