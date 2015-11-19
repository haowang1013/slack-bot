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

func pingHandler(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error {
	utils.SendMessage(rtm, strings.Join(fields, " "), m.Channel)
	return nil
}

func execHandler(rtm *slack.RTM, fields []string, m *slack.MessageEvent) error {
	if len(fields) == 0 {
		return errors.New("Missing command name")
	}

	go func(rtm *slack.RTM, fields []string, m *slack.MessageEvent) {
		cmd := exec.Command(fields[0], fields[1:]...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			msg := fmt.Sprintf("Failed to run '%s': %s", strings.Join(fields, " "), err.Error())
			utils.Log.Error(msg)
			utils.SendMessage(rtm, msg, m.Channel)
		} else {
			utils.SendMessage(rtm, string(out), m.Channel)
		}
	}(rtm, fields, m)
	return nil
}

func init() {
	addHandler("ping", HandlerFunc(pingHandler))
	addHandler("exec", HandlerFunc(execHandler))
}

func HandleMessage(rtm *slack.RTM, m *slack.MessageEvent) {
	fields := strings.Fields(m.Text)
	cmd := strings.ToLower(fields[0])
	if handler, ok := handlers[cmd]; ok {
		begin := time.Now()
		err := handler.HandleCommand(rtm, fields[1:], m)
		if err == nil {
			duration := time.Now().Sub(begin)
			utils.Log.Info("Command '%s' executed in %.2f sec", cmd, duration.Seconds())
		} else {
			msg := fmt.Sprintf("Command '%s' failed: %s", cmd, err.Error())
			utils.Log.Error(msg)
			utils.SendMessage(rtm, msg, m.Channel)
		}
	} else {
		// command not recognized
	}
}
