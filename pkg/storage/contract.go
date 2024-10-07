package storage

import "context"

const DEFAULT_NEWS_COUNT = 10

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

	// Articles returns slice of *pageSize* articles starting from *pageNum*
	Articles(ctx context.Context, pageNum int, pageSize int) ([]Article, error)

	// SearchArticles returns *pageSize* articles starting from page *pageNum* with search string in their title
	SearchArticles(ctx context.Context, search string, pageNum int, pageSize int) ([]Article, error)

	// Feeds returns slice of all added RSS feeds
	Feeds(context.Context) ([]RSSFeed, error)

	// AddArticle создаёт новую статью в хранилище из переданной структуры Article
	AddArticle(context.Context, Article) (int, error)

	// AddArticles сохраняет переданный срез []Article
	AddArticles(context.Context, []Article) error

	// DeleteArticle удалает из хранилища переданный Article
	DeleteArticle(context.Context, Article) error

	// AddRSSFeed создаёт в хранилище новый источник RSS-лент
	AddRSSFeed(context.Context, RSSFeed) (int, error)

	// DeleteRSSFeed удаляет из хранилища переданный RSSFeed
	DeleteRSSFeed(context.Context, RSSFeed) error
}
