package api

import "go-news-scraper/pkg/storage"

type Paginator struct {
	NumPages int `json:"num_pages"`
	CurPage  int `json:"cur_page"`
	PageSize int `json:"page_size"`
}

type PaginatedResponse struct {
	Articles  []storage.ArticleShort `json:"articles"`
	Paginator Paginator              `json:"paginator"`
}
