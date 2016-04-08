package commands

import (
	"errors"
	"fmt"
	"github.com/haowang1013/slack-bot/utils"
	"github.com/nlopes/slack"
	"os/exec"
	"strings"
	"time"
)

var rtm *slack.RTM

func init() {
	rootHandler.addHandler("ping", HandlerFunc(pingHandler))
	rootHandler.addHandler("exec", HandlerFunc(execHandler))
}

func SetRTM(_rtm *slack.RTM) {
	rtm = _rtm
}

func HandleMessage(m *slack.MessageEvent) {
	fields := strings.Fields(m.Text)
	cmd := strings.ToLower(fields[0])

	if handler, ok := rootHandler[cmd]; ok {
		begin := time.Now()
		err := handler.HandleCommand(fields[1:], m)
		if err == nil {
			duration := time.Now().Sub(begin)
			utils.Log.Infof("Command '%s' executed in %.2f sec", cmd, duration.Seconds())
		} else {
			msg := fmt.Sprintf("Command '%s' failed: %s", cmd, err.Error())
			utils.Log.Error(msg)
			sendMessage(msg, m.Channel)
		}
	} else {
		// command not recognized
	}
}

func sendMessage(message string, channel string) {
	utils.SendMessage(rtm, message, channel)
}

func pingHandler(fields []string, m *slack.MessageEvent) error {
	sendMessage(strings.Join(fields, " "), m.Channel)
	return nil
}

func execHandler(fields []string, m *slack.MessageEvent) error {
	if len(fields) == 0 {
		return errors.New("Missing command name")
	}

	go func(fields []string, m *slack.MessageEvent) {
		cmd := exec.Command(fields[0], fields[1:]...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			msg := fmt.Sprintf("Failed to run '%s': %s", strings.Join(fields, " "), err.Error())
			utils.Log.Error(msg)
			sendMessage(msg, m.Channel)
		} else {
			sendMessage(string(out), m.Channel)
		}
	}(fields, m)
	return nil
}
