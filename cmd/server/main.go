package main

import (
	"context"
	"go-news-scraper/pkg/api"
	"go-news-scraper/pkg/config"
	"go-news-scraper/pkg/scraper"
	"go-news-scraper/pkg/storage"
	"go-news-scraper/pkg/storage/pg"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	AppName = "go-news-scraper"
)

type server struct {
	db  storage.Storage
	api *api.API
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if cfg == nil {
		os.Exit(0)
	}

	log.Printf("starting %s service\n", AppName)
	log.Println(config.VersionString())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		cancel()
	}()

	db, err := pg.New(cfg.DBConnString)
	if err != nil {
		log.Fatal(err)
	}

	db.LoadRSSFeeds(ctx, cfg)

	rssScraper := scraper.New(db, cfg.RequestPeriod)
	rssScraper.Run(ctx)

	server := api.New(cfg.HttpPort, db)
	if err := server.Run(ctx); err != nil {
		log.Println("Got error:", err)
		os.Exit(0)
	}
}
