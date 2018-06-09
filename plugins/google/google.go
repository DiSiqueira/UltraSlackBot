package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/disiqueira/ultraslackbot/pkg/bot"
	"github.com/disiqueira/ultraslackbot/pkg/plugin"
	"github.com/disiqueira/ultraslackbot/pkg/slack"
)

const (
	pattern          = "(?i)\\b(google|search|gse)\\b"
	googleKeyEnvName = "GOOGLEKEY"
	googleCXEnvName  = "GOOGLECX"
	searchURL        = "https://www.googleapis.com/customsearch/v1"
)

type (
	google struct {
		plugin.BasicCommand
		cx        string
		googleKey string
	}

	googleResponse struct {
		Items []struct {
			Snippet string `json:"snippet"`
			Title   string `json:"title"`
			Link    string `json:"link"`
		} `json:"items"`
	}
)

func (gi *google) Start(specs bot.Specs) error {
	key, ok := specs.Get(googleKeyEnvName)
	if !ok {
		return fmt.Errorf("config %s not found", googleKeyEnvName)
	}
	gi.googleKey = key.(string)

	cx, ok := specs.Get(googleCXEnvName)
	if !ok {
		return fmt.Errorf("config %s not found", googleCXEnvName)
	}
	gi.cx = cx.(string)

	return nil
}

func (gi *google) Name() string {
	return "google"
}

func (gi *google) Execute(event slack.Event, botUser bot.UserInfo) ([]slack.Message, error) {
	return gi.HandleEvent(event, botUser, gi.matcher, gi.command)
}

func (gi *google) matcher() *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func (gi *google) command(text string) (string, error) {
	args := strings.Split(strings.TrimSpace(text), " ")
	if len(args) < 2 {
		return "", nil
	}

	text = strings.Join(args[1:], " ")

	var gisURL *url.URL
	gisURL, err := url.Parse(searchURL)
	if err != nil {
		return "", err
	}

	parameters := url.Values{}
	parameters.Add("key", gi.googleKey)
	parameters.Add("cx", gi.cx)
	parameters.Add("q", text)
	gisURL.RawQuery = parameters.Encode()

	data := &googleResponse{}
	if err := plugin.GetJSON(gisURL.String(), data); err != nil {
		return "", err
	}

	link := ""
	for _, item := range data.Items {
		if len(item.Link) > 0 {
			link = item.Link
			break
		}
	}

	return link, nil
}

var CustomPlugin google