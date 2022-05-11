/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"os"

	"github.com/govirtuo/kubectl-suspender/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List namespaces managed by kube-ns-suspender and their state",
	Long: `The list command lists all the namespaces that are managed by 
kube-ns-suspender. It also shows their current state.

Usage example:
	kubectl suspender list`,
	Run: func(cmd *cobra.Command, args []string) {
		s := utils.CreateSpinner("listing namespaces...")
		s.Start()

		// list all the namespaces
		ns, err := a.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			a.Logger.Fatal().Err(err).Msg("cannot list namespaces")
		}

		// create the table that will hold the results
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "State"})

		// populate the namespaces in the data variable
		var data [][]string
		for _, n := range ns.Items {
			val, ok := n.Annotations[a.AnnotationsNames.Prefix+"/"+a.AnnotationsNames.ControllerName]
			if !ok {
				continue
			}
			if val == a.ControllerName {
				state, ok := n.Annotations[a.AnnotationsNames.Prefix+"/"+a.AnnotationsNames.DesiredState]
				if !ok {
					continue
				}
				data = append(data, []string{
					n.Name, state,
				})
			}
		}

		for _, v := range data {
			table.Append(v)
		}

		s.Stop()
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
