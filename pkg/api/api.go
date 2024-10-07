package api

import (
	"context"
	"fmt"
	"go-news-scraper/pkg/storage"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	HttpListenPort int
	db             storage.Storage
	mux            *mux.Router
}

func New(port int, db storage.Storage) *API {
	api := API{
		HttpListenPort: port,
		db:             db,
		mux:            mux.NewRouter(),
	}

	return &api
}

func (api *API) Run(ctx context.Context) error {
	errChan := make(chan error)
	srv := api.serve(ctx, errChan)

	select {
	case <-ctx.Done():
		log.Println("gracefully shutting down")
		srv.Shutdown(ctx)
		return ctx.Err()
	case err := <-errChan:
		log.Println(err)
		return err
	}
}

func (api *API) serve(ctx context.Context, errChan chan error) *http.Server {
	api.initMux()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", api.HttpListenPort),
		Handler: api.mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, fmt.Sprintf(":%v", api.HttpListenPort), l.Addr().String())
			return ctx
		},
	}

	go func(s *http.Server) {
		if err := s.ListenAndServe(); err != nil {
			errChan <- err
		}
	}(httpServer)

	log.Println("serving HTTP server at", api.HttpListenPort)

	return httpServer
}

func (api *API) initMux() {
	api.mux.HandleFunc("/news", api.getNews).Methods(http.MethodGet, http.MethodOptions)
	api.mux.HandleFunc("/news/{id}", api.getArticle).Methods(http.MethodGet, http.MethodOptions)
	api.mux.HandleFunc("/rss", api.getFeeds).Methods(http.MethodGet, http.MethodOptions)
	api.mux.HandleFunc("/rss", api.addRSSFeed).Methods(http.MethodPost, http.MethodOptions)
	api.mux.HandleFunc("/rss", api.deleteRSSFeed).Methods(http.MethodDelete, http.MethodOptions)
	api.mux.Use(URLSchemaMiddleware(api.mux))
	api.mux.Use(RequestIDLoggerMiddleware(api.mux))
	api.mux.Use(LoggerMiddleware(api.mux))
}
