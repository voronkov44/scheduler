package main

import (
	"flag"
	"log/slog"
	"os"

	"scheduler/config"
	"scheduler/internal/calendar"
	"scheduler/internal/platform/postgres"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "config file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	db, err := postgres.New(cfg.DB)
	if err != nil {
		log.Error("cannot open db", "error", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(
		&calendar.Event{},
	); err != nil {
		log.Error("automigrate failed", "error", err)
		os.Exit(1)
	}

	log.Info("automigrate done")
}
