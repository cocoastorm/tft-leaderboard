package web

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Options struct {
	ServeAddress      string        `mapstructure:"address"`
	ServeAppDirPath   string        `mapstructure:"app-path"`
	ServeAppIndexPath string        `mapstructure:"app-indexpath"`
	ServeWriteTimeout time.Duration `mapstructure:"write-timeout"`
	ServeReadTimeout  time.Duration `mapstructure:"read-timeout"`
	Poll              bool          `mapstructure:"poll"`
	PollInterval      time.Duration `mapstructure:"poll-interval"`
	GoalRank          string        `mapstructure:"goal-rank"`
}

var DefaultOptions = Options{
	ServeAddress:      ":8080",
	ServeAppDirPath:   "static",
	ServeAppIndexPath: "index.html",
	ServeWriteTimeout: time.Second * 30,
	ServeReadTimeout:  time.Second * 30,
	Poll:              true,
	PollInterval:      time.Minute * 2,
	GoalRank:          "MASTERS",
}

type WebApp struct {
	Board *Board
	Sync  *RankSync
}

func (w *WebApp) Run(ctx context.Context, opts *Options) error {
	srv, err := w.Board.BuildServer(opts)
	if err != nil {
		return err
	}

	// start goroutine for: graceful web server shutdown (?)
	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		srv.Shutdown(shutdownCtx)
	}()

	// start polling
	if opts.Poll {
		go w.Sync.Do(ctx)
	}

	// start web server
	return srv.ListenAndServe()
}

// frontHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type frontHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the front handler. If a file is found, it will be served. If not, the
// file located at the index path on the front handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h frontHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
