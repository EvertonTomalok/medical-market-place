package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "go-template",
	Short: "Backend Go lang Template",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	_ = viper.BindEnv("App.Host", "APP_HOST")
	_ = viper.BindEnv("App.Port", "APP_PORT")
	_ = viper.BindEnv("App.Database.Host", "APP_DATABASE_HOST")
	_ = viper.BindEnv("App.Database.Port", "APP_DATABASE_PORT")
	_ = viper.BindEnv("App.Database.Name", "APP_DATABASE_NAME")
	_ = viper.BindEnv("App.Database.ConnMaxLifetime", "APP_DATABASE_CONN_MAX_LIFE_TIME")
	_ = viper.BindEnv("App.Database.MaxOpenConnections", "APP_DATABASE_MAX_OPEN_CONNECTIONS")
	_ = viper.BindEnv("App.Database.MaxIdleConnections", "APP_DATABASE_MAX_IDLE_CONNECTIONS")
}
