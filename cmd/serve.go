package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"ktn-x.com/tft-leaderboard/web"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Short: "Boots up the web server",
	Long: `Boots up the web server and static SPA routing
		for the leaderboard portion of tft-leaderboard
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// web server options from cmd
		opts := webOptions()

		// access Riot Games API
		riotCli := openRiotApi()

		// access data storage
		storeDB, err := openDB()
		if err != nil {
			log.Fatalf("failed to open DB: %s", err)
		}

		// web server + polling ranks
		app := web.WebApp{
			Board: web.NewBoard(storeDB),
			Sync: web.NewPoll(
				riotCli,
				storeDB,
				opts.PollInterval,
			),
		}

		// handle things gracefully
		// os.Interrupt -> ctrl-c
		// SIGTERM -> docker/k8s
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		// 1, 2, 3.
		// Blast off!
		if err = app.Run(ctx, opts); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	},
}

var (
	serveAddress string
	serveAppDirPath string
	serveAppIndexPath string
	serveWriteTimeout time.Duration
	serveReadTimeout time.Duration
	pollInterval time.Duration
	goalRank string
)

func init() {
	serveCmd.Flags().StringVar(&serveAddress, "address", web.DefaultOptions.ServeAddress, "address for web server to listen to")
	serveCmd.Flags().StringVar(&serveAppDirPath, "app-path", web.DefaultOptions.ServeAppDirPath, "path to built front app")
	serveCmd.Flags().StringVar(&serveAppIndexPath, "app-indexpath", web.DefaultOptions.ServeAppIndexPath, "index file to serve from built front app")
	serveCmd.Flags().DurationVar(&serveWriteTimeout, "write-timeout", web.DefaultOptions.ServeWriteTimeout, "server write timeout")
	serveCmd.Flags().DurationVar(&serveReadTimeout, "read-timeout", web.DefaultOptions.ServeReadTimeout, "server read timeout")
	serveCmd.Flags().DurationVar(&pollInterval, "poll", web.DefaultOptions.PollInterval, "how many minutes to poll/update for ranked data")
	serveCmd.Flags().StringVar(&goalRank, "goal", web.DefaultOptions.GoalRank, "the rank goal (eg. MASTERS)")

	viper.BindPFlag("address", serveCmd.Flags().Lookup("address"))
	viper.BindPFlag("app-path", serveCmd.Flags().Lookup("app-path"))
	viper.BindPFlag("app-indexpath", serveCmd.Flags().Lookup("app-indexpath"))
	viper.BindPFlag("write-timeout", serveCmd.Flags().Lookup("write-timeout"))
	viper.BindPFlag("read-timeout", serveCmd.Flags().Lookup("read-timeout"))
	viper.BindPFlag("poll", serveCmd.Flags().Lookup("poll"))
	viper.BindPFlag("goal", serveCmd.Flags().Lookup("goal"))

	rootCmd.AddCommand(serveCmd)
}

func webOptions() *web.Options {
	return &web.Options{
		ServeAddress: serveAddress,
		ServeAppDirPath: serveAppDirPath,
		ServeAppIndexPath: serveAppIndexPath,
		ServeWriteTimeout: serveWriteTimeout,
		ServeReadTimeout: serveReadTimeout,
		PollInterval: pollInterval,
		GoalRank: goalRank,
	}
}
