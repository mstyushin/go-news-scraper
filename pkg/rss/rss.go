// Пакет для работы с RSS-потоками.
// аккуратно украдено из примера реализации
package rss

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/mstyushin/go-news-scraper/pkg/storage"

	strip "github.com/grokify/html-strip-tags-go"
)

type FeedInfo struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

// Parse читает rss-поток и возвращет
// массив раскодированных новостей.
func Parse(f FeedInfo, feedID int) ([]storage.Article, error) {
	var data []storage.Article

	for _, item := range f.Channel.Items {
		var a storage.Article
		a.Title = item.Title
		a.Content = item.Description
		a.Content = strip.StripTags(a.Content)
		a.Link = item.Link
		a.RSSFeedID = feedID
		// Sat, 15 May 2021 04:05:00 +0300
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			a.PubTime = t.Unix()
		}
		data = append(data, a)
	}

	return data, nil
}
