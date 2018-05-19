package main

import (
	"strings"
	"github.com/disiqueira/ultraslackbot/pkg/slack"
	"github.com/disiqueira/ultraslackbot/pkg/bot"
	"github.com/CrowdSurge/banner"
)

type (
	echo struct { }
)

func (c *echo) Name() string {
	return "echo"
}

func (c *echo) Execute(event slack.Event, botUser bot.UserInfo) ([]slack.Message, error) {
	msg, err := slack.EventToMessage(event)
	if err != nil {
		return nil, nil
	}
	return c.handleMessageEvent(msg, botUser)
}

func (c *echo) handleMessageEvent(messageEvent slack.Message, botUser bot.UserInfo) ([]slack.Message, error) {
	args := strings.Split(strings.TrimSpace(messageEvent.Text()), " ")
	if len(args) < 2 || args[0] != c.Name() {
		return nil, nil
	}

	text := strings.Join(args[1:], " ")
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)

	outMessages := []slack.Message{
		slack.NewMessage("```"+banner.PrintS(text)+"```", messageEvent.Channel(), botUser),
		}

	return outMessages, nil
}

var CustomPlugin echo
