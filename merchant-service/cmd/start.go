package cmd

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the HTTP Server",
	Run: func(cmd *cobra.Command, args []string) {
		// app.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
