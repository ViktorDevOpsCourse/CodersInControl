//Copyright © 2023 Viktor Zhabskyi project on devops course viktordevopscourse@gmail.com

package cmd

import (
	"fmt"
	"github.com/viktordevopscourse/codersincontrol/app/internal/config"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.AppVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
