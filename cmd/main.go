package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/antonbaks/otusProjectJob/internal/app"
	"github.com/antonbaks/otusProjectJob/internal/cleaner"
	"github.com/antonbaks/otusProjectJob/internal/downloader"
	"github.com/antonbaks/otusProjectJob/internal/lru"
	"github.com/antonbaks/otusProjectJob/internal/resizer"
	"github.com/antonbaks/otusProjectJob/internal/server"
	"github.com/antonbaks/otusProjectJob/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/app/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	cfg := NewConfig(configFile)

	clearCh := make(chan string, 100)
	defer func() {
		close(clearCh)
		cancel()
	}()

	cache := lru.NewCache(cfg.CacheSize, clearCh)
	newStorage := storage.NewStorage(cfg.UploadDir)
	newDownloader := downloader.NewDownloader(&http.Client{})
	newResizer := resizer.NewResizer(cfg.MaxWidth, cfg.MinWidth, cfg.MaxHeight, cfg.MinHeight)

	newApp := app.NewApp(cache, newStorage, newDownloader, newResizer)

	httpServer := server.NewServer(&newApp, cfg.HTTPHost, cfg.HTTPPort)
	newCleaner := cleaner.NewCleaner(ctx, clearCh, newStorage)

	go func() {
		<-ctx.Done()

		if err := httpServer.Stop(); err != nil {
			log.Fatalln("Error stop http server: ", err)
		}
	}()

	go func() {
		newCleaner.Start()
	}()

	if err := httpServer.Start(); err != nil {
		log.Fatalln("Error start http server: ", err) // nolint:gocritic
	}
}
