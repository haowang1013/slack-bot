package commands

import (
	"errors"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"strings"
)

// if http handling should be async
var asyncHttp = false

var httpHandlers = handleCollection{
	"get": HandlerFunc(getHandler),
}

func init() {
	rootHandler.addHandler("http", HandlerFunc(httpHandler))
}

func getHandler(fields []string, m *slack.MessageEvent) error {
	resp, err := utils.HttpGet(fields[0])
	if err != nil {
		return err
	}
	sendMessage(string(resp), m.Channel)
	return nil
}

func fixUrl(url string) string {
	// note that url will be passed in as <url> by slack
	url = strings.TrimLeft(url, "<")
	url = strings.TrimRight(url, ">")
	return url
}

func httpHandler(fields []string, m *slack.MessageEvent) error {
	if len(fields) != 2 {
		return errors.New("Expecting params: <method> <url>")
	}

	// assuming fields[1] is the url
	fields[1] = fixUrl(fields[1])

	if asyncHttp {
		go httpHandlers.HandleCommand(fields, m)
		return nil
	} else {
		return httpHandlers.HandleCommand(fields, m)
	}
}
