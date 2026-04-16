package cmd

import (
	"fmt"
	"os"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Web API server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		mode := viper.GetString("mode")
		addr := fmt.Sprintf(":%s", port)

		fmt.Printf("Starting API server in %s mode on %s...\n", mode, addr)
		srv := api.NewServer(mode)
		if err := srv.Start(addr); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
}
