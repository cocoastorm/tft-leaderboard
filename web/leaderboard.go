package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

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
	api.HandleFunc("/leaderboard", l.Index)

	// spa front handler (eg. React App)
	router.PathPrefix("/").Handler(frontHandler{
		staticPath: opts.ServeAppDirPath,
		indexPath: opts.ServeAppIndexPath,
	})

	return router
}

func (l *Board) Index(w http.ResponseWriter, r *http.Request) {
	results, err := l.Store.ListContestantRanks()
	sort.Sort(results)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to fetch contestants and their ranks"))
		log.Printf("failed to index contestants & ranks: %s", err)
		return
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
