package main

import (
	"context"
	"flag"
	"time"

	"github.com/hazcod/notion2sen/config"
	"github.com/hazcod/notion2sen/pkg/notion"
	msSentinel "github.com/hazcod/notion2sen/pkg/sentinel"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	confFile := flag.String("config", "config.yml", "The YAML configuration file.")
	flag.Parse()

	cfg := config.Config{}
	if err := cfg.Load(*confFile); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("failed to load configuration")
	}

	if err := cfg.Validate(); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("invalid configuration")
	}

	logrusLevel, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.SetLevel(logrusLevel)

	//

	lookbackDuration, err := time.ParseDuration(cfg.Notion.Lookback)
	if err != nil {
		logger.WithError(err).WithField("lookback", cfg.Notion.Lookback).Fatal("invalid lookback duration provided")
	}
	lookback := time.Now().Add(-lookbackDuration)

	notionClient, err := notion.New(logger, cfg.Notion.ApiToken, cfg.Notion.OrganisationID)
	if err != nil {
		logger.WithError(err).Fatal("could not create notion client")
	}

	sentinel, err := msSentinel.New(logger, msSentinel.Credentials{
		TenantID:       cfg.Microsoft.TenantID,
		ClientID:       cfg.Microsoft.AppID,
		ClientSecret:   cfg.Microsoft.SecretKey,
		SubscriptionID: cfg.Microsoft.SubscriptionID,
		ResourceGroup:  cfg.Microsoft.ResourceGroup,
		WorkspaceName:  cfg.Microsoft.WorkspaceName,
	})
	if err != nil {
		logger.WithError(err).Fatal("could not create MS Sentinel client")
	}

	//

	if cfg.Microsoft.UpdateTable {
		if err := sentinel.CreateTable(ctx, logger, cfg.Microsoft.RetentionDays); err != nil {
			logger.WithError(err).Fatal("failed to create MS Sentinel table")
		}
	}

	//

	logger.WithField("lookback", lookback).Info("Retrieving Notion logs")

	logs, err := notionClient.GetAuditLogs(lookback)
	if err != nil {
		logger.WithError(err).Fatal("could not fetch onepassword signin events")
	}

	siemLogs, err := notionClient.ConvertAuditLogsToMap(logs)
	if err != nil {
		logger.WithError(err).Errorf("could not parse signin events")
	}

	//

	if err := sentinel.SendLogs(ctx, logger,
		cfg.Microsoft.DataCollection.Endpoint,
		cfg.Microsoft.DataCollection.RuleID,
		cfg.Microsoft.DataCollection.StreamName,
		siemLogs); err != nil {
		logger.WithError(err).Fatal("could not ship logs to sentinel")
	}

	//

	logger.WithField("total", len(siemLogs)).Info("successfully sent logs to sentinel")
}
