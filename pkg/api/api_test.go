package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/mstyushin/go-news-scraper/pkg/storage"
	"github.com/mstyushin/go-news-scraper/pkg/storage/pg"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

const dbURL = "postgres://postgres@localhost:5432/news?sslmode=disable"

var (
	s             *pg.DB
	pool          *pgxpool.Pool
	api           *API
	ctx           context.Context
	testRequestID = "f1548771-a62a-4068-8e28-73dca9a20a89"
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	var err error
	s, err = pg.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pool, err = pgxpool.Connect(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	api = New(8080, s)
	api.initMux()

	os.Exit(m.Run())
}

func TestAPI_getArticle(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/news/1?request_id=%s", testRequestID), nil)
	rr := httptest.NewRecorder()

	api.mux.ServeHTTP(rr, req)
	assert.True(t, rr.Code == http.StatusOK)
	assert.Equal(t, testRequestID, rr.Header().Get("x-request-id"), "should populate x-request-id header")

	b, err := ioutil.ReadAll(rr.Body)
	assert.NoError(t, err, "should be able to read response body")

	var data storage.Article
	err = json.Unmarshal(b, &data)
	assert.NoError(t, err, "response should be serializable")
	assert.Equal(t, 1, data.ID, "trying to get article with ID=1")
}

func TestAPI_getNews(t *testing.T) {
	articleTitle := "api test 1"
	articleContent := "something very thoughtful and interesting"
	articleRSSFeedID := 1
	pubTime := int64(time.Now().Unix())
	link := "https://news.ycombinator.com/item?id=39891148"

	a := storage.Article{
		RSSFeedID: articleRSSFeedID,
		Title:     articleTitle,
		Content:   articleContent,
		Link:      link,
		PubTime:   pubTime,
	}

	aid, _ := api.db.AddArticle(ctx, a)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/news?request_id=%s", testRequestID), nil)
	rr := httptest.NewRecorder()

	api.mux.ServeHTTP(rr, req)
	assert.True(t, rr.Code == http.StatusOK)
	assert.Equal(t, testRequestID, rr.Header().Get("x-request-id"), "should populate x-request-id header")

	b, err := ioutil.ReadAll(rr.Body)
	assert.NoError(t, err, "should be able to read response body")

	var data PaginatedResponse
	err = json.Unmarshal(b, &data)
	assert.NoError(t, err, "response should be serializable")
	assert.Equal(t, articleTitle, data.Articles[0].Title, "reading Title of a previously added article")

	a.ID = aid
	api.db.DeleteArticle(ctx, a)
}
