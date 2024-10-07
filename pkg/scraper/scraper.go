package scraper

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"go-news-scraper/pkg/rss"
	"go-news-scraper/pkg/storage"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Scraper отвечает за выборку статей из новостных лент и запись
// их в хранилище
type Scraper struct {
	storage        storage.Storage
	scrapeInterval time.Duration
	chErrs         chan error
	chArticles     chan []storage.Article
}

func New(s storage.Storage, scrapeInterval time.Duration) *Scraper {
	return &Scraper{
		storage:        s,
		scrapeInterval: scrapeInterval,
		chErrs:         make(chan error),
		chArticles:     make(chan []storage.Article),
	}
}

// Run запускает все горутины Scraper'a
func (s *Scraper) Run(ctx context.Context) {
	log.Println("starting scraper")

	// главная горутина Scraper'a которая занимается выборкой и парсингом
	go s.runUnwrapped(ctx)

	// горутина для перехвата ошибок
	go func() {
		for err := range s.chErrs {
			log.Println("scraper error:", err)
		}
	}()

	// горутина для сохранения статей
	go func() {
		for articles := range s.chArticles {
			err := s.storage.AddArticles(ctx, articles)
			if err != nil {
				log.Println("adding articles from feed")
				s.chErrs <- err
				continue
			}
		}
	}()
}

func (s *Scraper) runUnwrapped(ctx context.Context) {
	ticker := time.NewTicker(s.scrapeInterval)
	defer ticker.Stop()

	log.Println("scraping feeds")
	err := s.scrape(ctx)

	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				fmt.Println("scraper stopped")
				return
			}
		case <-ticker.C:
			log.Println("scraping feeds")
			err = s.scrape(ctx)
			if err != nil {
				s.chErrs <- err
				continue
			}
		}
	}
}

func (s *Scraper) scrape(ctx context.Context) error {
	log.Println("reading feeds from DB")
	feeds, err := s.storage.Feeds(ctx)
	if err != nil {
		return err
	}
	for _, feedObj := range feeds {

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedObj.URL, nil)
		if err != nil {
			log.Printf("client: could not create request: %s\n", err)
			return err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("client: error making http request: %s\n", err)
			return err
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var fi rss.FeedInfo
		err = xml.Unmarshal(b, &fi)
		if err != nil {
			log.Println("client: unmarshalling xml feed")
			return err
		}
		articles, err := rss.Parse(fi, feedObj.ID)
		if err != nil {
			log.Println("parsing rss feed")
			return err
		}

		s.chArticles <- articles
	}

	return nil
}
