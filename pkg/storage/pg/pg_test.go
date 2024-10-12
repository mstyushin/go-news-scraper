package pg

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mstyushin/go-news-scraper/pkg/config"
	"github.com/mstyushin/go-news-scraper/pkg/storage"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

const dbURL = "postgres://postgres@localhost:5432/news?sslmode=disable"

var s *DB
var pool *pgxpool.Pool
var ctx context.Context

func TestMain(m *testing.M) {
	var err error
	s, err = New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}

	ctx = context.Background()
	os.Exit(m.Run())
}

func getArticleByID(id int) (storage.Article, error) {
	ctx := context.Background()
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
	var res []storage.Article
	if err := pgxscan.Select(ctx, pool, &res, sql, id); err != nil {
		return storage.Article{}, err
	}

	if len(res) == 1 {
		return res[0], nil
	}

	return storage.Article{}, errors.New("unexpected search result")
}

func TestPG_AddArticle(t *testing.T) {
	articleTitle := "test"
	articleContent := "something very thoughtful and interesting"
	articleRSSFeedID := 1
	pubTime := int64(time.Now().Unix())

	a := storage.Article{
		RSSFeedID: articleRSSFeedID,
		Title:     articleTitle,
		Content:   articleContent,
		PubTime:   pubTime,
	}

	articleID, err := s.AddArticle(ctx, a)
	assert.NoError(t, err, "adding article")

	t.Log("Created article: ", articleID)

	_article, err := getArticleByID(articleID)
	assert.NoError(t, err, "searching for article")
	assert.Equal(t, articleTitle, _article.Title)
	assert.Equal(t, articleContent, _article.Content)
	assert.Equal(t, pubTime, _article.PubTime)
}

func TestPG_DeleteArticle(t *testing.T) {
	articleID := 1

	article := storage.Article{
		ID: articleID,
	}

	err := s.DeleteArticle(ctx, article)
	assert.NoError(t, err, "deleting article")

	t.Log("deleted article: ", articleID)

	_, err = getArticleByID(articleID)
	assert.EqualError(t, err, "unexpected search result")
}

func TestPG_SearchArticle(t *testing.T) {
	articleTitle := "Find me baby"
	articleContent := "whaterver you might be interested at"
	articleRSSFeedID := 1
	link := "https://news.ycombinator.com/item?id=39891948"

	pubTime := int64(time.Now().Unix())

	a := storage.Article{
		RSSFeedID: articleRSSFeedID,
		Title:     articleTitle,
		Content:   articleContent,
		PubTime:   pubTime,
		Link:      link,
	}

	articleID, err := s.AddArticle(ctx, a)
	assert.NoError(t, err, "adding article")

	t.Log("Created article: ", articleID)

	found, err := s.SearchArticles(ctx, "baby", 1, 10)
	assert.NoError(t, err, "searching articles")

	assert.NotEmpty(t, found, "should find at least one")
	assert.Equal(t, pubTime, found[0].PubTime)
}

func TestPG_LoadRSSFeeds(t *testing.T) {
	c := config.DefaultConfig()
	err := s.LoadRSSFeeds(ctx, c)
	assert.NoError(t, err, "loading RSS feeds from default config")
}
