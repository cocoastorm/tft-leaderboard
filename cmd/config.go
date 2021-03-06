package cmd

import (
	"fmt"
	"strings"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"ktn-x.com/tft-leaderboard/data"
	"ktn-x.com/tft-leaderboard/tft"
)

var cfgFile string

func bindGlobalConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tft.yml)")
	cmd.PersistentFlags().String("db", "tft-leaderboard.db", "Database Path")
	cmd.PersistentFlags().String("api-key", "", "Riot API Key")
	cmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	viper.BindPFlag("db", cmd.PersistentFlags().Lookup("db"))
	viper.BindPFlag("api-key", cmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("use-viper", cmd.PersistentFlags().Lookup("viper"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// config file settings
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".tft")

		// env settings
		viper.SetEnvPrefix("tft")

		// allow '-' and '_'
		replacer := strings.NewReplacer("-", "_")
		viper.SetEnvKeyReplacer(replacer)

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// config file not found, ignoring error
				// fallback to env
				viper.AutomaticEnv()
				return
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

func openDB() (*data.Store, error) {
	storagePath := viper.GetString("db")
	return data.OpenDB(storagePath)
}

func openRiotApi() *tft.RiotClient {
	apiKey := viper.GetString("api-key")
	return tft.NewRiot(apiKey)
}
