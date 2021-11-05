package main

import (
	"os"

	"github.com/grafana/flagger-k6-webhook/pkg"
	"github.com/grafana/flagger-k6-webhook/pkg/k6"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
)

const (
	flagCloudToken = "cloud-token"
	flagLogLevel   = "log-level"
	flagListenPort = "listen-port"
	flagSlackToken = "slack-token"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("execution failed: %s", err)
	}
}

func run(args []string) error {
	app := cli.NewApp()
	app.Name = "flagger-k6-webhook"
	app.Usage = "Launches k6 load testing from a flagger webhook"
	app.Action = launchServer

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    flagCloudToken,
			EnvVars: []string{"K6_CLOUD_TOKEN"},
		},
		&cli.IntFlag{
			Name:    flagListenPort,
			EnvVars: []string{"LISTEN_PORT"},
			Value:   80,
		},
		&cli.StringFlag{
			Name:    flagLogLevel,
			EnvVars: []string{"LOG_LEVEL"},
			Value:   log.InfoLevel.String(),
		},
		&cli.StringFlag{
			Name:    flagSlackToken,
			EnvVars: []string{"SLACK_TOKEN"},
		},
	}

	return app.Run(args)
}

func launchServer(c *cli.Context) error {
	logLevel, err := log.ParseLevel(c.String(flagLogLevel))
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)

	client, err := k6.NewClient(c.String(flagCloudToken))
	if err != nil {
		return err
	}

	var slackClient *slack.Client
	if slackToken := c.String(flagSlackToken); slackToken != "" {
		slackClient = slack.New(slackToken)
	}

	return pkg.Listen(client, slackClient, c.Int(flagListenPort))
}
