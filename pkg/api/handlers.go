package api

import (
	"encoding/json"
	"go-news-scraper/pkg/storage"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (api *API) getNews(w http.ResponseWriter, r *http.Request) {
	var s string
	var articles []storage.Article
	var err error
	// TODO make it configurable
	pageSize := 10

	pageNum := 1

	if r.URL.Query().Has("page_size") {
		s = r.URL.Query().Get("page_size")
		pageSize, _ = strconv.Atoi(s)
	}

	if r.URL.Query().Has("page") {
		s = r.URL.Query().Get("page")
		pageNum, _ = strconv.Atoi(s)
	}

	if r.URL.Query().Has("s") {
		searchString := r.URL.Query().Get("s")
		articles, err = api.db.SearchArticles(r.Context(), searchString, pageNum, pageSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {

		articles, err = api.db.Articles(r.Context(), pageNum, pageSize)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	paginator := Paginator{
		NumPages: len(articles),
		CurPage:  pageNum,
		PageSize: pageSize,
	}
	paginated := &PaginatedResponse{
		Articles:  articles,
		Paginator: paginator,
	}

	bytes, err := json.Marshal(paginated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (api *API) getArticle(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	aid, _ := strconv.Atoi(s)
	article, err := api.db.Article(r.Context(), aid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func (api *API) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := api.db.Feeds(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bytes, err := json.Marshal(feeds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}

func (api *API) addRSSFeed(w http.ResponseWriter, r *http.Request) {
	var feed storage.RSSFeed
	err := json.NewDecoder(r.Body).Decode(&feed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = api.db.AddRSSFeed(r.Context(), feed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) deleteRSSFeed(w http.ResponseWriter, r *http.Request) {
	var feed storage.RSSFeed
	err := json.NewDecoder(r.Body).Decode(&feed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.DeleteRSSFeed(r.Context(), feed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) deleteArticle(w http.ResponseWriter, r *http.Request) {
	var a storage.Article
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.DeleteArticle(r.Context(), a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
