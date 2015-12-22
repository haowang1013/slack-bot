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
		"log":  HandlerFunc(getMessagesHandler),
		"send": HandlerFunc(sendMessagesHandler),
	})
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
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	temp := make([]map[string]interface{}, 0)
	err = json.Unmarshal(content, &temp)
	if err != nil {
		return err
	}

	marshalled, err := json.MarshalIndent(temp, "", " ")
	if err != nil {
		return err
	}
	sendMessage(string(marshalled), m.Channel)
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

	utils.Log.Debug(string(buff))

	req, err := http.NewRequest("POST", "https://leancloud.cn/1.1/rtm/messages", strings.NewReader(string(buff)))
	if err != nil {
		return err
	}
	req.Header.Add("X-LC-Id", appID)
	req.Header.Add("X-LC-Key", fmt.Sprintf("%s,master", masterKey))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
