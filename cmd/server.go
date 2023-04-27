package cmd

import (
	"github.com/EvertonTomalok/marketplace-health/internal/app"
	"github.com/EvertonTomalok/marketplace-health/internal/ui/rest"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run http server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New()
		ctx := cmd.Context()
		defer app.CloseConnections(ctx)

		config := app.Configure(ctx)

		switch config.App.LogLevel {
		case "ERROR", "error":
			logger.SetLevel(log.ErrorLevel)
		case "INFO", "info":
			logger.SetLevel(log.InfoLevel)
		case "DEBUG", "debug":
			logger.SetLevel(log.DebugLevel)
		default:
			logger.SetLevel(log.ErrorLevel)
		}

		app.InitDB(ctx, &config, logger)
		rest.RunServer(ctx, config)
	},
}

func init() {
	serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(serverCmd)
}
