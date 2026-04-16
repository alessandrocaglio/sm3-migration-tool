package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alessandrocaglio/sm3-migration-tool/pkg/checkers"
	"github.com/alessandrocaglio/sm3-migration-tool/pkg/discovery"
	"github.com/alessandrocaglio/sm3-migration-tool/pkg/report"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan the cluster for Service Mesh 2 resources and check compatibility",
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
			loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
			configOverrides := &clientcmd.ConfigOverrides{}
			kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
			
			restConfig, err := kubeConfig.ClientConfig()
			if err != nil {
				fmt.Printf("Error loading kubeconfig: %v\n", err)
				os.Exit(1)
			}

			disc, err = discovery.NewLiveDiscovery(config, restConfig)
			if err != nil {
				fmt.Printf("Error creating live discovery: %v\n", err)
				os.Exit(1)
			}
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
			checkers.NewAddonsChecker(),
			checkers.NewNetworkPolicyChecker(),
			checkers.NewRouteChecker(),
		}

		var allFindings []checkers.Finding
		for _, c := range allCheckers {
			findings, err := c.Check(context.Background(), state)
			if err != nil {
				fmt.Printf("Error in checker %s: %v\n", c.Name(), err)
				continue
			}
			allFindings = append(allFindings, findings...)
		}

		fmt.Print(report.GenerateMarkdown(allFindings))
		
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile != "" {
			err := report.GenerateJSON(allFindings, outputFile)
			if err != nil {
				fmt.Printf("Error generating JSON report: %v\n", err)
			} else {
				fmt.Printf("\nJSON report saved to %s\n", outputFile)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("output", "o", "", "Output file for JSON report")
}
