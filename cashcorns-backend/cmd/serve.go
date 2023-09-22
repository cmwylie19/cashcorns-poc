/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/cmwylie19/cashcorns-backend/pkg/server"
	"github.com/spf13/cobra"
)

var (
	// Port is the port the server will run on
	Port string = "443"

	// PayRunFileLocation is the location of the payrun file
	PayRunFileLocation string = "/etc/config/payrun.json"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run webserver",
	Long: `Run webserver with specific configuration.

Usage:
  cashcorns-backend serve

Flags:
--port, -p  Port to run the server on. Default is 443.
--payrun-file-location, -f  Location of the payrun file.`,
	Run: func(cmd *cobra.Command, args []string) {
		server := server.NewServer(Port, PayRunFileLocation)
		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVarP(&Port, "port", "p", "8888", "Port to run the server on. Default is 8888.")
	serveCmd.PersistentFlags().StringVarP(&PayRunFileLocation, "payrun-file-location", "f", "/etc/config/payrun.json", "Location of the payrun file.")
}
