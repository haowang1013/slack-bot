package commands

import (
	"errors"
	"fmt"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
	"strings"
)

// if http handling should be async
var asyncHttp = false

var httpHandlers = map[string]CommandHandler{}

func init() {
	httpHandlers["get"] = HandlerFunc(getHandler)
}

func getHandler(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error {
	url := fields[0]
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	utils.SendMessage(rtm, string(body), m.Channel)
	return nil
}

func fixUrl(url string) string {
	// note that url will be passed in as <url> by slack
	url = strings.TrimLeft(url, "<")
	url = strings.TrimRight(url, ">")
	return url
}

func httpHandler(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error {
	if len(fields) != 2 {
		return errors.New("Expecting params: <method> <url>")
	}

	method := strings.ToLower(fields[0])
	fields = fields[1:]
	fields[0] = fixUrl(fields[0])

	if h, ok := httpHandlers[method]; ok {
		if asyncHttp {
			go h.HandleCommand(rtm, fields, m)
			return nil
		} else {
			return h.HandleCommand(rtm, fields, m)
		}
	} else {
		return errors.New(fmt.Sprintf("Unsupported method: %s", method))
	}
}

func init() {
	addHandler("http", HandlerFunc(httpHandler))
}
