package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	appCommit = "n/a"
	appBuilt  = "n/a"
)

// @title ASCII Canvas Editor
// @version 1.0
// @description ASCII Canvas Editor
// @license.name Public Domain
// @BasePath /api

func main() {
	var rootCmd = &cobra.Command{
		Use:   "ascanvas",
		Short: "ASCII Canvas Editor",
		Run:   Serve,
	}

	var cmdServe = &cobra.Command{
		Use:   "serve",
		Short: "Start the web API server",
		Run:   Serve,
	}
	rootCmd.AddCommand(cmdServe)

	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Check software version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nCommit : %v\nBuilt: %v\n", appCommit, appBuilt)
		},
	}
	rootCmd.AddCommand(cmdVersion)

	var err = rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

type Config struct {
	ListenAddr string `json:"listen_addr"`
	DbDriver   string `json:"db_driver"`
	DSN        string `json:"dsn"`
	LogLevel   string `json:"log_level"`
}
