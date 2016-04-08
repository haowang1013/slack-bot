package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"os"
	"strconv"
)

const (
	envAppID     = "LEANCLOUD_APPID"
	envAppKey    = "LEANCLOUD_APPKEY"
	envMasterKey = "LEANCLOUD_MASTERKEY"
	chatUrl      = "https://api.leancloud.cn/1.1/classes/_Conversation"
	messageUrl   = "https://leancloud.cn/1.1/rtm/messages"
)

var (
	appID          string
	appKey         string
	masterKey      string
	defaultHeaders map[string]string
)

type messageBody struct {
	Sender    string `json:"from_peer"`
	Text      string `json:"message"`
	ChanID    string `json:"conv_id"`
	Transient bool   `json:"transient"`
}

type channel struct {
	Name     string `json:"name"`
	ObjectID string `json:"objectId"`
}

type channelQuery struct {
	Channels []channel `json:"results"`
}

type createChannelRequest struct {
	Name string `json:"name"`
}

func init() {
	found := false
	appID, found = os.LookupEnv(envAppID)
	if !found {
		utils.Log.Errorf("Missing environment variable '%s'", envAppID)
	}

	appKey, found = os.LookupEnv(envAppKey)
	if !found {
		utils.Log.Errorf("Missing environment variable '%s'", envAppKey)
	}

	masterKey, found = os.LookupEnv(envMasterKey)
	if !found {
		utils.Log.Errorf("Missing environment variable '%s'", envMasterKey)
	}

	defaultHeaders = map[string]string{
		"X-LC-Id":  appID,
		"X-LC-Key": appKey,
	}

	rootHandler.addHandler("chat", &handleCollection{
		"channels":   HandlerFunc(listChannelsHandler),
		"log":        HandlerFunc(getMessagesHandler),
		"send":       HandlerFunc(sendMessagesHandler),
		"bulkcreate": HandlerFunc(bulkCreateChannelsHandler),
		"find":       HandlerFunc(findChannelsHandler),
	})
}

func prettyJson(obj interface{}) (string, error) {
	content, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		return "", err
	}
	return string(content), nil
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

func findChannel(name string) (*channelQuery, error) {
	query := &channelQuery{}
	url := fmt.Sprintf("%s?where={\"name\":\"%s\"}", chatUrl, name)
	err := utils.NewJsonRequest("GET", url, &defaultHeaders, nil, query)
	if err != nil {
		return nil, err
	}
	return query, nil
}

func findChannelsHandler(fields []string, m *slack.MessageEvent) error {
	if err := checkEnv(); err != nil {
		return err
	}

	if len(fields) != 1 {
		return errors.New("Expecting params: <name>")
	}
	query, err := findChannel(fields[0])
	if err != nil {
		return err
	}
	content, err := prettyJson(query)
	if err != nil {
		return err
	}
	sendMessage(utils.FormatMessage(content, "```"), m.Channel)
	return nil
}

func listChannelsHandler(fields []string, m *slack.MessageEvent) error {
	if err := checkEnv(); err != nil {
		return err
	}

	response := make(map[string]interface{})
	err := utils.NewJsonRequest("GET", chatUrl, &defaultHeaders, nil, &response)
	if err != nil {
		return err
	}
	content, err := prettyJson(response)
	if err != nil {
		return err
	}
	sendMessage(utils.FormatMessage(content, "```"), m.Channel)
	return nil
}

func bulkCreateChannelsHandler(fields []string, m *slack.MessageEvent) error {
	if err := checkEnv(); err != nil {
		return err
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
		channelNames = append(channelNames, fmt.Sprintf("%s:Channel_%d", tenant, i+1))
	}

	numCreated := 0
	numExisting := 0
	for _, name := range channelNames {
		query, err := findChannel(name)
		if err != nil {
			return err
		}

		numChannels := len(query.Channels)
		if numChannels > 0 {
			utils.Log.Debugf("Found %d channels under name '%s'", numChannels, name)
			numExisting++
		} else {
			r := createChannelRequest{
				Name: name,
			}

			err := utils.NewJsonRequest("POST", chatUrl, &defaultHeaders, r, nil)
			if err != nil {
				return err
			}
			utils.Log.Debugf("Channel '%s' created", name)
			numCreated++
		}
	}
	sendMessage(fmt.Sprintf("Created %d channels, %d already exist", numCreated, numExisting), m.Channel)
	return nil
}

func getMessagesHandler(fields []string, m *slack.MessageEvent) error {
	if err := checkEnv(); err != nil {
		return err
	}

	if len(fields) == 0 {
		return errors.New("Missing channel ID")
	}

	uri := fmt.Sprintf("%s/logs?convid=%s&limit=10", messageUrl, fields[0])
	response := make(map[string]interface{})
	err := utils.NewJsonRequest("GET", uri, &defaultHeaders, nil, &response)
	if err != nil {
		return err
	}

	content, err := prettyJson(response)
	if err != nil {
		return err
	}
	sendMessage(utils.FormatMessage(content, "```"), m.Channel)
	return nil
}

func sendMessagesHandler(fields []string, m *slack.MessageEvent) error {
	if err := checkEnv(); err != nil {
		return err
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

	headers := map[string]string{
		"X-LC-Id":  appID,
		"X-LC-Key": fmt.Sprintf("%s,master", masterKey),
	}
	err := utils.NewJsonRequest("POST", messageUrl, &headers, msg, nil)
	if err != nil {
		return err
	}

	sendMessage("message sent", m.Channel)
	return nil
}
