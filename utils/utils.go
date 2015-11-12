package utils

import (
	"github.com/kr/pretty"
	"github.com/nlopes/slack"
)

func ListGroups(api *slack.Client) {
	groups, err := api.GetGroups(false)
	if err != nil {
		return
	}

	for _, group := range groups {
		Log.Debug("%v", pretty.Formatter(group))
	}
}

func ListChannels(api *slack.Client) {
	channels, err := api.GetChannels(false)
	if err != nil {
		return
	}

	for _, c := range channels {
		Log.Debug("Channel '%s': '%s'\n", c.Name, c.ID)
	}
}

func SendMessage(rtm *slack.RTM, message string, channel string) {
	rtm.SendMessage(rtm.NewOutgoingMessage(message, channel))
}
