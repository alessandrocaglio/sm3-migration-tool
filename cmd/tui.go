package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/checkers"
	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
	"github.com/alessandrocaglio/sm3-migration-tool/pkg/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI migration assistant",
	Run: func(cmd *cobra.Command, args []string) {
		mode := viper.GetString("mode")
		namespace := viper.GetString("namespace")
		var disc discovery.Discovery

		config := discovery.DiscoveryConfig{
			ControlPlaneNamespace: namespace,
		}

		if mode == "mock" {
			disc = discovery.NewMockDiscovery("testdata/mock-cluster.yaml", config)
		} else {
			fmt.Println("Live discovery not yet implemented. Use --mode mock")
			os.Exit(1)
		}

		state, err := disc.Discover(context.Background())
		if err != nil {
			fmt.Printf("Error during discovery: %v\n", err)
			os.Exit(1)
		}

		allCheckers := []checkers.Checker{
			checkers.NewGatewayChecker(),
			checkers.NewServiceEntryChecker(),
			checkers.NewMtlsChecker(),
		}

		var allFindings []checkers.Finding
		for _, c := range allCheckers {
			findings, err := c.Check(context.Background(), state)
			if err != nil {
				continue
			}
			allFindings = append(allFindings, findings...)
		}

		if err := tui.Run(allFindings); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
