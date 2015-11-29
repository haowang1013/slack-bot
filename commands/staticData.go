package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
)

const (
	indexUrl    = "https://s3-eu-west-1.amazonaws.com/directive-tiers.dg-api.com/static-data/directivegames/the-machines-static-data/index.json"
	dataRootUrl = "https://d1pqpvg8ar5xgy.cloudfront.net/directivegames/the-machines-static-data/data/"
)

type staticDataEntry struct {
	CommitID string `json:"commit_id"`
	Ref      string `json:"ref"`
	Types    string `json:"types"`
	Schemas  string `json:"schemas"`
}

type staticDataRoot struct {
	Entries []staticDataEntry `json:"index"`
}

func init() {
	rootHandler.addHandler("sd", &handleCollection{
		"commits": HandlerFunc(commitsHandler),
		"commit":  HandlerFunc(commitHandler),
	})
}

func commitsHandler(fields []string, m *slack.MessageEvent) error {
	resp, err := utils.HttpGet(indexUrl)
	if err != nil {
		return err
	}

	var root staticDataRoot
	err = json.Unmarshal(resp, &root)
	if err != nil {
		return err
	}

	for i := range root.Entries {
		elem := &root.Entries[i]
		elem.Types, elem.Schemas = getDataUrl(elem.CommitID)
		buff, err := json.MarshalIndent(elem, "", " ")
		if err != nil {
			return err
		}
		sendMessage(string(buff), m.Channel)
	}
	return nil
}

func getDataUrl(commitID string) (types string, schemas string) {
	types = fmt.Sprintf("%s%s/types.json", dataRootUrl, commitID)
	schemas = fmt.Sprintf("%s%s/schemas.json", dataRootUrl, commitID)
	return
}

func commitHandler(fields []string, m *slack.MessageEvent) error {
	if len(fields) == 0 {
		return errors.New("Missing commit id")
	}

	types, schemas := getDataUrl(fields[0])
	msg := fmt.Sprintf("types: %s\nschemas: %s", types, schemas)
	sendMessage(msg, m.Channel)
	return nil
}
