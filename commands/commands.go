package commands

import (
	//"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"github.com/nlopes/slack"
)

func HandleMessage(m *slack.MessageEvent) {
	if len(m.BotID) > 0 {
		// ignore bot message
		return
	}

	fmt.Printf("Got message: %v\n", pretty.Formatter(m))
}
