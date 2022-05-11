package cmd

import (
	"context"

	"github.com/govirtuo/kubectl-suspender/app"
	"github.com/govirtuo/kubectl-suspender/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// unsuspendCmd represents the unsus command
var unsuspendCmd = &cobra.Command{
	Use:   "unsuspend",
	Short: "Unsuspend a namespace managed by kube-ns-suspender",
	Long: `The unsuspend command allows a user to unsuspend a namespace
managed by kube-ns-suspender. If no namespace is provided in arg,
it will prompt a list of namespaces to the user.

Usage example:
	kubectl suspender unsuspend
	kubectl suspender unsuspend my-namespace
	kubectl suspender unsuspend my-namespace my-other-namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if args are provided
		// if yes: consider the args as being namespaces names
		// if no: prompt a list of namespaces and let the user
		//   decide which one we should suspend
		var namespacesToSuspend []string
		switch len(args) {
		case 0:
			s := utils.CreateSpinner("listing namespaces...")
			s.Start()
			ns, err := a.Clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				a.Logger.Fatal().Err(err).Msg("cannot list namespaces")
			}

			// get the list of namespaces that are watched by kube-ns-suspender and are running
			cnannot := a.AnnotationsNames.Prefix + "/" + a.AnnotationsNames.ControllerName
			dsannot := a.AnnotationsNames.Prefix + "/" + a.AnnotationsNames.DesiredState
			wl := utils.GetWatchedAndSuspendedNamespaces(ns, cnannot, dsannot, a.ControllerName)

			if len(wl) == 0 {
				s.Stop()
				a.Logger.Info().Msg("no namespace is currently suspended")
				return
			}

			// create the prompt and populate the item list with the running namespaces
			var l []string
			for _, n := range wl {
				l = append(l, n.Name)
			}
			s.Stop()
			prompt := promptui.Select{
				Label: "Select a namespace to unsuspend",
				Items: l,
			}

			_, result, err := prompt.Run()
			if err != nil {
				a.Logger.Fatal().Err(err).Msg("cannot prompt namespaces list")
			}
			namespacesToSuspend = append(namespacesToSuspend, result)
		default:
			// handle the args as being namespaces names to suspend
			namespacesToSuspend = args
		}

		// suspend the namespaces here based on what is in namespacesToSuspend
		a.Logger.Info().Msgf("list of namespaces that will be unsuspended: %s", namespacesToSuspend)

		for _, n := range namespacesToSuspend {
			status := "done\n"
			s := utils.CreateSpinner("unsuspending namespace " + n)
			s.Start()
			if err := a.UpdateNamespace(n, app.RunningState); err != nil {
				a.Logger.Error().Err(err).Msgf("cannot unsuspend namespace '%s'", n)
				status = "failed\n"
			}
			s.FinalMSG = "unsuspending namespace " + n + ": " + status
			s.Stop()
		}
	},
}

func init() {
	rootCmd.AddCommand(unsuspendCmd)
}
