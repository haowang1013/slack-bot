package utils

import (
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"github.com/nlopes/slack"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	client = &http.Client{}
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

func FormatMessage(msg string, format string) string {
	return fmt.Sprintf("%s%s%s", format, msg, format)
}

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func NewJsonRequest(verb string, url string, headers *map[string]string, payload interface{}, response interface{}) error {
	// prepare the payload
	var body io.Reader
	if payload != nil {
		buff, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		body = strings.NewReader(string(buff))
	}

	// make the request
	req, err := http.NewRequest(verb, url, body)
	if err != nil {
		return err
	}

	// add the headers
	if headers != nil {
		for k, v := range *headers {
			req.Header.Add(k, v)
		}
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// parse the response
	if response != nil {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(content, response)
		if err != nil {
			Log.Error("Failed to unmarshall json string '%s' from request '%s'", string(content), url)
			return err
		}
	}

	return nil
}
