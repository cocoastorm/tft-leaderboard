package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"ktn-x.com/tft-leaderboard/app"
)

func wrapImportError(err error) error {
	return fmt.Errorf("failed to import: %s", err)
}

var (
	filename string

	importCmd = &cobra.Command{
		Use: "import",
		Short: "Imports participants into the leaderboard",
		Long: "Imports tft summoner participants into the tft leaderboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			// access Riot Games API
			riotCli := openRiotApi()

			// access data storage
			storeDB, err := openDB()
			if err != nil {
				return fmt.Errorf("failed to open DB: %s", err)
			}

			defer storeDB.Close()

			// jabber app
			jabber := app.Jabber{
				Riot: riotCli,
				Store: storeDB,
			}

			// participants pool
			pool := app.NewParticipantsFile(filename)

			// pool: error with file
			if err := pool.Error(); err != nil {
				return wrapImportError(err)
			}
			
			contestants, err := jabber.Import(pool)
			if err != nil {
				return wrapImportError(err)
			}

			// pool: error with decoding
			if err := pool.Error(); err != nil {
				return wrapImportError(err)
			}

			// list contestants for debug purposes
			for _, c := range contestants {
				summoner := c.Summoner
				fmt.Printf("[%d] %s - id: %s\n", c.SequenceId, summoner.Name, summoner.Id)
			}

			return nil
		},
	}
)

func init() {
	importCmd.Flags().StringVar(&filename, "i", "", "input file to import tft participants in")
	rootCmd.AddCommand(importCmd)
}
