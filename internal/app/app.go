package app

import (
	"context"
	"github.com/disiqueira/ultraslackbot/pkg/command/news"
	"log"
	"os"

	"github.com/disiqueira/ultraslackbot/pkg/command/nba"
	"github.com/disiqueira/ultraslackbot/pkg/command/pokequest"
	"github.com/disiqueira/ultraslackbot/pkg/command/uptime"

	"github.com/spf13/cobra"

	"github.com/disiqueira/ultraslackbot/internal/bot"
	"github.com/disiqueira/ultraslackbot/internal/conf"
	usbCtx "github.com/disiqueira/ultraslackbot/internal/context"
	"github.com/disiqueira/ultraslackbot/internal/handlers/admin"
	"github.com/disiqueira/ultraslackbot/internal/handlers/command"
	"github.com/disiqueira/ultraslackbot/internal/handlers/logger"
	"github.com/disiqueira/ultraslackbot/pkg/command/9gag"
	"github.com/disiqueira/ultraslackbot/pkg/command/aftership"
	"github.com/disiqueira/ultraslackbot/pkg/command/bible"
	"github.com/disiqueira/ultraslackbot/pkg/command/cat"
	"github.com/disiqueira/ultraslackbot/pkg/command/choose"
	"github.com/disiqueira/ultraslackbot/pkg/command/chucknorris"
	"github.com/disiqueira/ultraslackbot/pkg/command/docs"
	"github.com/disiqueira/ultraslackbot/pkg/command/echo"
	"github.com/disiqueira/ultraslackbot/pkg/command/emoji"
	"github.com/disiqueira/ultraslackbot/pkg/command/fortune"
	"github.com/disiqueira/ultraslackbot/pkg/command/google"
	"github.com/disiqueira/ultraslackbot/pkg/command/hello"
	"github.com/disiqueira/ultraslackbot/pkg/command/howlongtobeat"
	"github.com/disiqueira/ultraslackbot/pkg/command/isup"
	"github.com/disiqueira/ultraslackbot/pkg/command/lastfm"
	"github.com/disiqueira/ultraslackbot/pkg/command/lenny"
	"github.com/disiqueira/ultraslackbot/pkg/command/random"
	"github.com/disiqueira/ultraslackbot/pkg/command/seedr"
	"github.com/disiqueira/ultraslackbot/pkg/command/shrug"
	"github.com/disiqueira/ultraslackbot/pkg/command/urban"
	"github.com/disiqueira/ultraslackbot/pkg/command/uuid"
	"github.com/disiqueira/ultraslackbot/pkg/command/watch"
	"github.com/disiqueira/ultraslackbot/pkg/command/wikipedia"
	"github.com/disiqueira/ultraslackbot/pkg/command/wolfram"
	"github.com/disiqueira/ultraslackbot/pkg/command/youtube"
	"github.com/disiqueira/ultraslackbot/pkg/slack"
)

type (
	App struct{}
)

const (
	successExitCode      = 0
	errorExitCode        = 1
	slackTokenEnvVar     = "SLACKTOKEN"
	googleKeyEnvVar      = "GOOGLEKEY"
	googleCXEnvVar       = "GOOGLECX"
	wolframKeyEnvName    = "WOLFRAMKEY"
	lastFMKeyEnvName     = "LASTFMKEY"
	afterShipKeyEnvName  = "AFTERSHIPKEY"
	seedrUsernameEnvName = "SEEDRUSERNAME"
	seedrPasswordEnvName = "SEEDRPASSWORD"
	newsAPIKeyEnvName    = "NEWSAPIKEY"
)

func New() *App {
	return &App{}
}

func (a *App) Run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	outLogger := log.New(os.Stdout, "", 0)
	errLogger := log.New(os.Stderr, "", 0)

	ctx = usbCtx.WithOutLogger(ctx, outLogger)
	ctx = usbCtx.WithErrLogger(ctx, errLogger)

	specs := conf.Load()

	outLogger.Printf("Configurations: %+v", specs.All())

	slackToken, ok := specs.Get(slackTokenEnvVar)
	if !ok {
		errLogger.Printf("config missing %s", slackTokenEnvVar)
		os.Exit(errorExitCode)
	}

	key, ok := specs.Get(googleKeyEnvVar)
	if !ok {
		errLogger.Printf("config missing %s", googleKeyEnvVar)
		os.Exit(errorExitCode)
	}

	cx, ok := specs.Get(googleCXEnvVar)
	if !ok {
		errLogger.Printf("config missing %s", googleCXEnvVar)
		os.Exit(errorExitCode)
	}

	wolframKey, ok := specs.Get(wolframKeyEnvName)
	if !ok {
		errLogger.Printf("config missing %s", wolframKeyEnvName)
		os.Exit(errorExitCode)
	}

	lastfmKey, ok := specs.Get(lastFMKeyEnvName)
	if !ok {
		errLogger.Printf("config missing %s", lastFMKeyEnvName)
		os.Exit(errorExitCode)
	}

	afterShipKey, ok := specs.Get(afterShipKeyEnvName)
	if !ok {
		errLogger.Printf("config missing %s", afterShipKeyEnvName)
		os.Exit(errorExitCode)
	}

	seedrUsername, ok := specs.Get(seedrUsernameEnvName)
	if !ok {
		errLogger.Printf("config missing %s", seedrUsernameEnvName)
		os.Exit(errorExitCode)
	}

	seedrPassword, ok := specs.Get(seedrPasswordEnvName)
	if !ok {
		errLogger.Printf("config missing %s", seedrPasswordEnvName)
		os.Exit(errorExitCode)
	}

	newsAPIKey, ok := specs.Get(newsAPIKeyEnvName)
	if !ok {
		errLogger.Printf("config missing %s", newsAPIKeyEnvName)
		os.Exit(errorExitCode)
	}

	slackClient := slack.New(slackToken.(string))
	b := bot.New(slackClient)

	commandList := []*cobra.Command{
		google.NewGoogleImageCommand(key.(string), cx.(string)),
		google.NewGoogleSearchCommand(key.(string), cx.(string)),
		hello.NewHelloCommand(),
		ninegag.New9gagCommand(),
		choose.NewChooseCommand(),
		youtube.NewYoutubeCommand(key.(string)),
		wolfram.NewWolframCommand(wolframKey.(string)),
		wikipedia.NewWikipediaCommand(),
		urban.NewUrbanCommand(),
		cat.NewCatCommand(),
		fortune.NewFortuneCommand(),
		howlongtobeat.NewHLTBCommand(),
		lenny.NewLennyCommand(),
		shrug.NewShrugCommand(),
		lastfm.NewLastFMCommand(lastfmKey.(string)),
		random.NewRandomCommand(),
		chucknorris.NewChuckNorrisFactCommand(),
		echo.NewEchoCommand(),
		emoji.NewEmojiCommand(),
		isup.NewIsUpCommand(),
		afterShip.NewAfterShipCommand(afterShipKey.(string)),
		uuid.NewUUIDCommand(),
		docs.NewDocsCommand(),
		bible.NewBibleCommand(),
		seedr.NewSeedrCommand(seedrUsername.(string), seedrPassword.(string)),
		watch.NewWatchCommand(seedrUsername.(string), seedrPassword.(string)),
		pokequest.NewPokeQuestCommand(),
		uptime.NewUptimeCommand(),
		nba.NewNBACommand(),
		news.NewNewsCommand(newsAPIKey.(string)),
	}

	handlerList := []bot.Handler{
		logger.New(),
		admin.New(b),
		command.New(commandList),
	}

	b.SetHandlers(handlerList)
	b.Run(ctx)
	os.Exit(successExitCode)
}
