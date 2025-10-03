package app

import (
	"context"
	"os"

	"github.com/JesterSe7en/scrapego/config"
	"github.com/JesterSe7en/scrapego/internal/logger"
	"github.com/JesterSe7en/scrapego/internal/presenter"
	"github.com/JesterSe7en/scrapego/internal/scraper"
	"github.com/JesterSe7en/scrapego/internal/storage"

	wp "github.com/JesterSe7en/scrapego/internal/workerpool"
)

type App struct {
	log       *logger.Logger
	cfg       *config.Config
	presenter presenter.Presenter
	cache     storage.CacheStorage
}

func New(log *logger.Logger, cfg *config.Config) (*App, error) {
	manager := storage.GetCacheManager()

	cache, err := manager.GetCache()
	if err != nil {
		return nil, err
	}

	var p presenter.Presenter
	if cfg.Format == "json" {
		p = presenter.NewJSONPresenter(os.Stdout)
	} else {
		p = presenter.NewTextPresenter(os.Stdout)
	}

	return &App{
		log:       log,
		cfg:       cfg,
		presenter: p,
		cache:     cache,
	}, nil
}

func (a *App) Run(ctx context.Context, urls []string) error {
	// create scraper client
	scraperClient := scraper.New(a.cfg, a.log, a.cache)

	// startup worker pool
	pool := wp.New(a.cfg.Concurrency, a.cfg.Concurrency)
	pool.Run(ctx, a.cfg.Concurrency)

	// create jobs for worker pool
	go func() {
		defer pool.Close()
		for _, url := range urls {
			scrapeURL := url
			job := func() wp.Result {
				return scraperClient.ScrapeWithRetry(scrapeURL)
			}

			if err := pool.Submit(ctx, job); err != nil {
				a.log.Warn("Shutting down job submission: %v", err)
				return
			}
		}
	}()

	// Process results as they come in
	for res := range pool.Results() {
		if res.Err != nil {
			a.log.Error("Failed to get response: %s", res.Err.Error())
			continue
		}

		if err := a.presenter.Write(res.Value); err != nil {
			a.log.Error("Failed to write output: %v", err)
		}
	}

	return nil
}
