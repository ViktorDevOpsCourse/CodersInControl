/*
Copyright © 2024 devops@codersincontr
*/

package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/viktordevopscourse/codersincontrol/app/internal/app"
	"github.com/viktordevopscourse/codersincontrol/app/pkg/logger"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slackbot",
	Short: "slackbot is a slack bot that can automatically deploy applications\n",
	Long:  `slackbot is slack bot that can automatically deploy applications.`,
	Run:   run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// run represents the Run function for rootCmd
func run(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	log := logger.GetDefaultLogger()
	logger.ToContext(ctx, log)

	err := app.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
