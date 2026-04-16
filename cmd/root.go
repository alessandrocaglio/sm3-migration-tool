package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFile string
	logger  *zap.Logger
)

var rootCmd = &cobra.Command{
	Use:   "sm3-migration-tool",
	Short: "A tool to assist in migrating from OSSM 2.x to 3.x",
	Long: `sm3-migration-tool is a utility designed to help OpenShift administrators
migrate their Service Mesh installations from version 2.6.14 to 3.0.0.
It performs discovery, compatibility checks, and generates migration reports.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sm3-migration-tool.yaml)")
	rootCmd.PersistentFlags().StringP("mode", "m", "live", "Execution mode: live or mock")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace of the control plane to scan")
	rootCmd.PersistentFlags().StringP("allowed-namespaces-regex", "r", ".*", "Regular expression to restrict allowed control plane namespaces")
	
	viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode"))
	viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("allowed-namespaces-regex", rootCmd.PersistentFlags().Lookup("allowed-namespaces-regex"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sm3-migration-tool")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initLogger() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
}
