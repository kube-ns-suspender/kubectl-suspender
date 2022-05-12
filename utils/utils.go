package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/govirtuo/kubectl-suspender/app"
	v1 "k8s.io/api/core/v1"
)

func GetWatchedAndRunningNamespaces(ns *v1.NamespaceList, controllerNameFullAnnot, desiredStateFullAnnot, controllerName string) []v1.Namespace {
	var watchlist []v1.Namespace
	for _, n := range ns.Items {
		val, ok := n.Annotations[controllerNameFullAnnot]
		if !ok || val != controllerName {
			continue
		}
		if n.Annotations[desiredStateFullAnnot] == app.RunningState {
			watchlist = append(watchlist, n)
		}
	}
	return watchlist
}

func GetWatchedAndSuspendedNamespaces(ns *v1.NamespaceList, controllerNameFullAnnot, desiredStateFullAnnot, controllerName string) []v1.Namespace {
	var watchlist []v1.Namespace
	for _, n := range ns.Items {
		val, ok := n.Annotations[controllerNameFullAnnot]
		if !ok || val != controllerName {
			continue
		}
		if n.Annotations[desiredStateFullAnnot] == app.SuspendedState {
			watchlist = append(watchlist, n)
		}
	}
	return watchlist
}

func CreateSpinner(suffix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + suffix
	return s
}

func IsNamespaceWatched(ns *v1.Namespace, controllerNameFullAnnot, controllerName string) error {
	val, ok := ns.Annotations[controllerNameFullAnnot]
	if !ok {
		return fmt.Errorf("this namespace is not watched by kube-ns-suspender")
	}
	if val != controllerName {
		return fmt.Errorf("namespace's controller name does not match the configured controller name")
	}
	return nil
}
