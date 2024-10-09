package storage

import "context"

type Article struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Link      string `json:"link"`
	RSSFeedID int    `json:"rss_feed_id"`
	PubTime   int64  `json:"pub_time"`
}

type ArticleShort struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	ShortContent string `json:"short_content"`
	LinkToFull   string `json:"link_to_full"`
	PubTime      int64  `json:"pub_time"`
}

type RSSFeed struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type Storage interface {
	// Article return one article by given ID
	Article(context.Context, int) (Article, error)

	// Articles returns slice of *pageSize* short articles starting from *pageNum*
	Articles(ctx context.Context, pageNum int, pageSize int) ([]Article, error)

	// SearchArticles returns *pageSize* articles starting from page *pageNum* with search string in their title
	SearchArticles(ctx context.Context, search string, pageNum int, pageSize int) ([]Article, error)

	// CountArticles returns articles count, filtered by search string (if non-empty)
	CountArticles(ctx context.Context, search string) (int, error)

	// Feeds returns slice of all added RSS feeds
	Feeds(context.Context) ([]RSSFeed, error)

	// AddArticle creates new Article and returns id assigned by Storage
	AddArticle(context.Context, Article) (int, error)

	// AddArticles stores []Article slice
	AddArticles(context.Context, []Article) error

	// DeleteArticle deletes Article
	DeleteArticle(context.Context, Article) error

	// AddRSSFeed creates RSS feed source and returns id assigned by Storage
	AddRSSFeed(context.Context, RSSFeed) (int, error)

	// DeleteRSSFeed deletes RSSFeed
	DeleteRSSFeed(context.Context, RSSFeed) error
}
