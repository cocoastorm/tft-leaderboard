package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"ktn-x.com/tft-leaderboard/data"
)

type Board struct {
	Store *data.Store;
}

func NewBoard(store *data.Store) *Board {
	return &Board{
		Store: store,
	}
}

func (l *Board) buildRouter(opts *Options) *mux.Router {
	router := mux.NewRouter()

	// api routes
	api := router.PathPrefix("/api").Subrouter()
	
	// leaderboard index
	api.HandleFunc("/leaderboard", l.Index).Methods("GET")

	// spa front handler (eg. React App)
	router.PathPrefix("/").Handler(frontHandler{
		staticPath: opts.ServeAppDirPath,
		indexPath: opts.ServeAppIndexPath,
	})

	return router
}

func (l *Board) Index(w http.ResponseWriter, r *http.Request) {
	ts, err := l.Store.GetRankTimestamp()
	if err != nil {
		ts = 0
		log.Printf("failed to fetch rank timestamp: %s", err)
	}

	// cache validation
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching#cache_validation
	if ifModifiedSince := r.Header.Get(""); ifModifiedSince != "" {
		rank := time.Unix(int64(ts), 0)
		
		since, err := time.Parse(http.TimeFormat, ifModifiedSince)
		if err != nil {
			log.Printf("failed to parse 'If-Modified-Since' cache header: %s", err)
		}

		if since.After(rank) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	results, err := l.Store.ListContestantRanks()
	sort.Sort(results)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to fetch contestants and their ranks"))
		log.Printf("failed to index contestants & ranks: %s", err)
		return
	}

	// cache headers
	w.Header().Set("Cache-Control", "max-age=120, must-validate")

	if ts > 0 {
		tsf := time.Unix(int64(ts), 0)
		w.Header().Set("Last-Modified", tsf.Format(http.TimeFormat))
	}

	err = json.NewEncoder(w).Encode(&results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("uh oh"))
		log.Printf("failed to encode contestants & ranks: %s", err)
		return
	}
}

// todo: add ctx.Background() if we're also booting up other stuffs
func (l *Board) BuildServer(opts *Options) (*http.Server, error) {
	if opts == nil {
		return nil, fmt.Errorf("web server config options are required")
	}

	router := l.buildRouter(opts)

	return &http.Server{
		Handler: router,
		Addr: opts.ServeAddress,
		// Good practice: enforce timeouts for serves you create!
		ReadTimeout: opts.ServeReadTimeout,
		WriteTimeout: opts.ServeWriteTimeout,
	}, nil
}
