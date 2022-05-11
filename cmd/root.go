package cmd

import (
	"os"

	"github.com/govirtuo/kubectl-suspender/app"
	"github.com/spf13/cobra"
)

var a app.App

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl suspender",
	Short: "A kubectl plugin to (un)suspend namespaces managed by kube-ns-suspender",
	Long: `This plugin for kubectl allows a user to suspend or unsuspend namespaces
managed by kube-ns-suspender directly from the command-line.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ap *app.App) {
	a = *ap
	rootCmd.Use = "kubectl-suspender"
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
}
