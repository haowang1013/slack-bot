package commands

import (
	"errors"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
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
	sendMessage(string(body), m.Channel)
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
