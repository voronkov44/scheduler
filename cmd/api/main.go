package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"scheduler/config"
	"scheduler/internal/integrations/itmo"
)

func main() {
	var configPath string

	flag.StringVar(
		&configPath,
		"config",
		"config.yaml",
		"config file",
	)
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Error(
			"cannot load timezone",
			"error", err,
		)
		os.Exit(1)
	}

	client := itmo.NewClient(
		cfg.ITMO.AdapterURL,
		cfg.ITMO.APIToken,
		cfg.ITMO.Timeout,
	)

	service := itmo.NewService(
		client,
		location,
	)

	now := time.Now().In(location)

	events, err := service.Fetch(
		context.Background(),
		now,
		now.AddDate(0, 0, 14),
	)
	if err != nil {
		log.Error(
			"cannot fetch ITMO schedule",
			"error", err,
		)
		os.Exit(1)
	}

	log.Info(
		"ITMO schedule fetched",
		"events", len(events),
	)

	for _, event := range events {
		log.Info(
			"ITMO event",
			"id", event.ExternalID,
			"title", event.Title,
			"start", event.StartsAt,
			"end", event.EndsAt,
			"location", event.Location,
		)
	}
}
