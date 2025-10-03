package app

import (
	"context"
	"fmt"

	"github.com/JesterSe7en/scrapego/config"
	"github.com/JesterSe7en/scrapego/internal/logger"
	"github.com/JesterSe7en/scrapego/internal/scraper"
	"github.com/JesterSe7en/scrapego/internal/storage"

	wp "github.com/JesterSe7en/scrapego/internal/workerpool"
)

type App struct {
	log   *logger.Logger
	cfg   *config.Config
	cache storage.CacheStorage
}

func New(log *logger.Logger, cfg *config.Config) (*App, error) {
	manager := storage.GetCacheManager()

	cache, err := manager.GetCache()
	if err != nil {
		return nil, err
	}

	return &App{
		log:   log,
		cfg:   cfg,
		cache: cache,
	}, nil
}

func (a *App) Run(ctx context.Context, urls []string) error {

	// create scraper client
	scraperClient := scraper.New(a.cfg, a.log, a.cache)
	// loop through

	pool := wp.New(a.cfg.Concurrency, a.cfg.Concurrency)
	pool.Run(a.cfg.Concurrency)

	for _, url := range urls {
		pool.Submit(func() wp.Result {
			return scraperClient.ScrapeWithRetry(url)
		})
	}

	pool.Close()

	for res := range pool.Results() {
		if res.Err != nil {
			a.log.Error("Failed to get response: %s", res.Err.Error())
			continue
		}

		// put the results into stdout
		fmt.Println(res.Value)
	}
	return nil
}
