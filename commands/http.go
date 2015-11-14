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

func getHandler(rtm *slack.RTM, url string, m *slack.MessageEvent) error {
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

func httpHandler(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error {
	if len(fields) != 2 {
		return errors.New("Expecting params: <method> <url>")
	}

	method := strings.ToLower(fields[0])
	url := fields[1] // note that url will be passed in as <url> by slack
	url = strings.TrimLeft(url, "<")
	url = strings.TrimRight(url, ">")

	if method == "get" {
		return getHandler(rtm, url, m)
	} else {
		return errors.New(fmt.Sprintf("Unsupported method: %s", method))
	}
}

func init() {
	addHandler("http", HandlerFunc(httpHandler))
}
