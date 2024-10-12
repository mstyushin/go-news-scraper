package pg

import (
	"context"
	"errors"
	"log"

	"github.com/mstyushin/go-news-scraper/pkg/config"
	"github.com/mstyushin/go-news-scraper/pkg/storage"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ storage.Storage = &DB{}

type DB struct {
	pool *pgxpool.Pool
}

func New(url string) (*DB, error) {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	db := &DB{
		pool: pool,
	}

	return db, nil
}

func (db *DB) Article(ctx context.Context, aid int) (storage.Article, error) {
	sql := `
SELECT
  id,
  rss_feed_id,
  title,
  content,
  link,
  pub_time
FROM articles
WHERE id = $1;
`
	var a storage.Article

	row := db.pool.QueryRow(ctx, sql, aid)
	err := row.Scan(
		&a.ID,
		&a.RSSFeedID,
		&a.Title,
		&a.Content,
		&a.Link,
		&a.PubTime,
	)
	if err != nil {
		return storage.Article{}, err
	}

	return a, nil

}

func (db *DB) Articles(ctx context.Context, pageNum, pageSize int) ([]storage.Article, error) {
	offset := (pageNum - 1) * pageSize

	sql := `
SELECT
  id,
  rss_feed_id,
  title,
  content,
  link,
  pub_time
FROM articles
ORDER BY pub_time DESC
LIMIT $1
OFFSET $2;
`
	articles := make([]storage.Article, 0, pageSize)

	rows, err := db.pool.Query(ctx, sql, pageSize, offset)
	for rows.Next() {
		var a storage.Article
		err = rows.Scan(
			&a.ID,
			&a.RSSFeedID,
			&a.Title,
			&a.Content,
			&a.Link,
			&a.PubTime,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (db *DB) SearchArticles(ctx context.Context, search string, pageNum, pageSize int) ([]storage.Article, error) {
	offset := (pageNum - 1) * pageSize

	sql := `
SELECT
  id,
  rss_feed_id,
  title,
  content,
  link,
  pub_time
FROM articles
WHERE title ILIKE $1
ORDER BY pub_time DESC
LIMIT $2
OFFSET $3;
`
	var articles []storage.Article

	rows, err := db.pool.Query(ctx, sql, "%"+search+"%", pageSize, offset)
	for rows.Next() {
		var a storage.Article
		err = rows.Scan(
			&a.ID,
			&a.RSSFeedID,
			&a.Title,
			&a.Content,
			&a.Link,
			&a.PubTime,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (db *DB) CountArticles(ctx context.Context, search string) (int, error) {
	sql := "SELECT count(*) FROM articles;"
	sqlFiltered := "SELECT count(*) FROM articles WHERE title ILIKE $1;"
	var count int

	if search != "" {
		rows, err := db.pool.Query(ctx, sqlFiltered, "%"+search+"%")
		if err != nil {
			return 0, err
		}
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				return 0, err
			}
		}
	} else {
		rows, err := db.pool.Query(ctx, sql)
		if err != nil {
			return 0, err
		}
		for rows.Next() {
			if err := rows.Scan(&count); err != nil {
				return 0, err
			}
		}
	}

	return count, nil
}

func (db *DB) AddArticle(ctx context.Context, a storage.Article) (int, error) {
	sql := `
INSERT INTO
  articles (
	rss_feed_id,
    title,
    content,
	link,
	pub_time
  )
VALUES
  (
    $1, $2, $3, $4, $5
  ) RETURNING id
`
	var id int
	err := db.pool.QueryRow(
		ctx,
		sql,
		a.RSSFeedID,
		a.Title,
		a.Content,
		a.Link,
		a.PubTime,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	log.Println("Created article", id)

	return id, nil
}

func (db *DB) AddArticles(ctx context.Context, articles []storage.Article) error {
	var err error
	for _, a := range articles {
		_, err = db.AddArticle(ctx, a)
	}
	// мы хотим попытаться добавить каждую статью и если хотя бы
	// одна не добавилась - возвращаем ошибку
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteArticle(ctx context.Context, a storage.Article) error {
	sql := "DELETE FROM articles WHERE id = $1"
	r, err := db.pool.Exec(ctx, sql, a.ID)
	if err != nil {
		return err
	}
	if r.RowsAffected() == 0 {
		log.Println("Article", a.ID, "not found")
		return errors.New("not found")
	}
	log.Println("Deleted article", a.ID)

	return nil
}

func (db *DB) Feeds(ctx context.Context) ([]storage.RSSFeed, error) {
	sql := `
SELECT
  id,
  url
FROM rss_feeds
ORDER BY id;
`
	var feeds []storage.RSSFeed
	if err := pgxscan.Select(ctx, db.pool, &feeds, sql); err != nil {
		return nil, err
	}

	return feeds, nil
}

func (db *DB) AddRSSFeed(ctx context.Context, feed storage.RSSFeed) (int, error) {
	sql := `
INSERT INTO
  rss_feeds (
	url
  )
VALUES
  (
    $1
  ) RETURNING id
`
	var id int
	err := db.pool.QueryRow(
		ctx,
		sql,
		feed.URL,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	log.Println("Added RSS feed", id, "url:", feed.URL)

	return id, nil
}

func (db *DB) DeleteRSSFeed(ctx context.Context, feed storage.RSSFeed) error {
	sql := "DELETE FROM rss_feeds WHERE id = $1"
	r, err := db.pool.Exec(ctx, sql, feed.ID)
	if err != nil {
		return err
	}
	if r.RowsAffected() == 0 {
		log.Println("RSS feed", feed.ID, "not found")
		return errors.New("not found")
	}
	log.Println("Deleted feed", feed.ID)

	return nil
}

func (db *DB) LoadRSSFeeds(ctx context.Context, cfg *config.Config) error {
	log.Println("Loading RSS feeds from config")

	for _, url := range cfg.RSSFeeds {
		_, err := db.AddRSSFeed(ctx, storage.RSSFeed{URL: url})
		if err != nil {
			return err
		}
	}

	return nil
}
