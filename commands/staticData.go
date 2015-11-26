package commands

import (
	"errors"
	"fmt"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
)

const (
	indexUrl    = "https://s3-eu-west-1.amazonaws.com/directive-tiers.dg-api.com/static-data/directivegames/the-machines-static-data/index.json"
	dataRootUrl = "https://d1pqpvg8ar5xgy.cloudfront.net/directivegames/the-machines-static-data/data/"
)

func init() {
	rootHandler.addHandler("sd", &handleCollection{
		"index":   HandlerFunc(indexHandler),
		"types":   HandlerFunc(typesHandler),
		"schemas": HandlerFunc(schemasHandler),
	})
}

func indexHandler(fields []string, m *slack.MessageEvent) error {
	resp, err := utils.HttpGet(indexUrl)
	if err != nil {
		return err
	}
	sendMessage(string(resp), m.Channel)
	return nil
}

func typesHandler(fields []string, m *slack.MessageEvent) error {
	if len(fields) == 0 {
		return errors.New("Missing commit id")
	}

	url := fmt.Sprintf("%s%s/types.json", dataRootUrl, fields[0])
	sendMessage(url, m.Channel)
	return nil
}

func schemasHandler(fields []string, m *slack.MessageEvent) error {
	if len(fields) == 0 {
		return errors.New("Missing commit id")
	}

	url := fmt.Sprintf("%s%s/schemas.json", dataRootUrl, fields[0])
	sendMessage(url, m.Channel)
	return nil
}
