package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	envAppID     = "LEANCLOUD_APPID"
	envAppKey    = "LEANCLOUD_APPKEY"
	envMasterKey = "LEANCLOUD_MASTERKEY"
)

var (
	appID     string
	appKey    string
	masterKey string
	client    = &http.Client{}
)

type messageBody struct {
	Sender    string `json:"from_peer"`
	Text      string `json:"message"`
	ChanID    string `json:"conv_id"`
	Transient bool   `json:"transient"`
}

type conversation struct {
	Name string `json:"name"`
}

type conversationQuery struct {
	Results []conversation `json:"results"`
}

type createChannelRequest struct {
	Name string `json:"name"`
}

func init() {
	found := false
	appID, found = os.LookupEnv(envAppID)
	if !found {
		utils.Log.Error("Missing environment variable '%s'", envAppID)
	}

	appKey, found = os.LookupEnv(envAppKey)
	if !found {
		utils.Log.Error("Missing environment variable '%s'", envAppKey)
	}

	masterKey, found = os.LookupEnv(envMasterKey)
	if !found {
		utils.Log.Error("Missing environment variable '%s'", envMasterKey)
	}

	rootHandler.addHandler("chat", &handleCollection{
		"channels":   HandlerFunc(listChannelsHandler),
		"log":        HandlerFunc(getMessagesHandler),
		"send":       HandlerFunc(sendMessagesHandler),
		"bulkcreate": HandlerFunc(bulkCreateChannelsHandler),
	})
}

func prettyJson(source []byte, isArray bool) ([]byte, error) {
	var temp interface{}
	if isArray {
		temp = make([]map[string]interface{}, 0)

	} else {
		temp = make(map[string]interface{})
	}

	if err := json.Unmarshal(source, &temp); err != nil {
		return nil, err
	}

	return json.MarshalIndent(temp, "", " ")
}

func checkEnv() error {
	if len(appID) == 0 {
		return errors.New("app ID is not set")
	}

	if len(appKey) == 0 {
		return errors.New("app key is not set")
	}

	if len(masterKey) == 0 {
		return errors.New("master key is not set")
	}

	return nil
}

func listChannelsHandler(fields []string, m *slack.MessageEvent) error {
	if err := checkEnv(); err != nil {
		return nil
	}

	req, err := http.NewRequest("GET", "https://leancloud.cn/1.1/classes/_Conversation", nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-LC-Id", appID)
	req.Header.Add("X-LC-Key", appKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	content, err = prettyJson(content, true)
	if err != nil {
		return err
	}

	sendMessage(string(content), m.Channel)
	return nil
}

func bulkCreateChannelsHandler(fields []string, m *slack.MessageEvent) error {
	if err := checkEnv(); err != nil {
		return nil
	}

	if len(fields) != 2 {
		return errors.New("Expecting params: <tenant> <count>")
	}

	tenant := fields[0]
	count, err := strconv.Atoi(fields[1])
	if err != nil {
		return err
	}

	channelNames := make([]string, 0)
	for i := 0; i < count; i++ {
		channelNames = append(channelNames, fmt.Sprintf("%s:Channel %d", tenant, i+1))
	}

	// get all the conversations
	req, err := http.NewRequest("GET", "https://api.leancloud.cn/1.1/classes/_Conversation", nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-LC-Id", appID)
	req.Header.Add("X-LC-Key", appKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	query := conversationQuery{}
	err = json.Unmarshal(content, &query)
	if err != nil {
		return err
	}

	// set of all existing channels
	names := make(map[string]bool)
	for _, c := range query.Results {
		names[c.Name] = true
	}

	// iterate all the ones we want and create missing ones
	numExisting := 0
	numCreated := 0
	for _, name := range channelNames {
		_, found := names[name]
		if found {
			utils.Log.Debug("Channel '%s' already exists", name)
			numExisting++
		} else {
			// not found, create a channel
			r := createChannelRequest{
				Name: name,
			}

			buff, err := json.Marshal(&r)
			if err != nil {
				return nil
			}

			req, err := http.NewRequest("POST", "https://api.leancloud.cn/1.1/classes/_Conversation", strings.NewReader(string(buff)))
			if err != nil {
				return err
			}
			req.Header.Add("X-LC-Id", appID)
			req.Header.Add("X-LC-Key", appKey)
			req.Header.Add("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			utils.Log.Debug("Channel '%s' created", name)
			numCreated++
		}
	}

	sendMessage(fmt.Sprintf("Created %d channels, %d already exist", numCreated, numExisting), m.Channel)
	return nil
}

func getMessagesHandler(fields []string, m *slack.MessageEvent) error {
	if err := checkEnv(); err != nil {
		return nil
	}

	if len(fields) == 0 {
		return errors.New("Missing channel ID")
	}

	uri := fmt.Sprintf("https://leancloud.cn/1.1/rtm/messages/logs?convid=%s&limit=10", fields[0])
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-LC-Id", appID)
	req.Header.Add("X-LC-Key", appKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	content, err = prettyJson(content, true)
	if err != nil {
		return err
	}

	sendMessage(string(content), m.Channel)
	return nil
}

func sendMessagesHandler(fields []string, m *slack.MessageEvent) error {
	if err := checkEnv(); err != nil {
		return nil
	}

	if len(fields) != 2 {
		return errors.New("Expecting params: <channel ID> <text>")
	}

	msg := messageBody{
		Sender:    "slack-bot",
		Text:      fields[1],
		ChanID:    fields[0],
		Transient: false,
	}

	buff, err := json.Marshal(&msg)
	if err != nil {
		return nil
	}

	req, err := http.NewRequest("POST", "https://leancloud.cn/1.1/rtm/messages", strings.NewReader(string(buff)))
	if err != nil {
		return err
	}
	req.Header.Add("X-LC-Id", appID)
	req.Header.Add("X-LC-Key", fmt.Sprintf("%s,master", masterKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sendMessage("message sent", m.Channel)
	return nil
}
