package api

import (
	"context"
	"encoding/json"
	"go-news-scraper/pkg/storage"
	"go-news-scraper/pkg/storage/pg"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

const dbURL = "postgres://postgres@localhost:5432/news?sslmode=disable"

var s *pg.DB
var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	s, err = pg.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestAPI_getArticle(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/news/1?request_id=12345", nil)
	rr := httptest.NewRecorder()

	api := New(8080, s)
	api.initMux()

	api.mux.ServeHTTP(rr, req)

	if !(rr.Code == http.StatusOK) {
		t.Errorf("wrong HTTP code: got %d, expected %d", rr.Code, http.StatusOK)
	}

	b, err := ioutil.ReadAll(rr.Body)
	assert.NoError(t, err, "should be able to read response body")

	var data storage.Article
	err = json.Unmarshal(b, &data)
	assert.NoError(t, err, "response should be serializable")
	assert.Equal(t, 1, data.ID, "trying to get article with ID=1")
}
