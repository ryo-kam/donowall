package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
)

func main() {
	appToken := os.Getenv("SLACK_APP_TOKEN")

	if appToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must be set.\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(appToken, "xapp-") {
		fmt.Fprintf(os.Stderr, "SLACK_APP_TOKEN must have the prefix \"xapp-\".")
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")

	if botToken == "" {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must be set.\n")
		os.Exit(1)
	}

	if !strings.HasPrefix(botToken, "xoxb-") {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN must have the prefix \"xoxb-\".")
	}

	slackApi := slack.New(
		botToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)

	slackClient := socketmode.New(
		slackApi,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	bttvClient, err := NewBttvClient()

	if err != nil {
		log.Fatal(err)
	}

	socketmodeHandler := socketmode.NewSocketmodeHandler(slackClient)

	socketmodeHandler.HandleSlashCommand("/donowall", donowallHandler(bttvClient))

	socketmodeHandler.RunEventLoop()
}
