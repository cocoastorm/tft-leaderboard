package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "tft-leaderboard",
	Short: "tft-leaderboard is a community tft rank leaderboard. race to diamond!",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("please use subcommand import or serve")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	bindGlobalConfigFlags(rootCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
