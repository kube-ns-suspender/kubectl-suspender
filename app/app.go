package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

const (
	RunningState   = "Running"
	SuspendedState = "Suspended"
)

type App struct {
	Logger           zerolog.Logger
	LogLevel         zerolog.Level
	LogLevelStr      string
	Clientset        *kubernetes.Clientset
	ControllerName   string
	AnnotationsNames Annotations
}

type Annotations struct {
	Prefix         string
	DesiredState   string
	ControllerName string
}

func New() (*App, error) {
	var a App
	var err error
	a.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	a.ControllerName = "kube-ns-suspender"
	a.AnnotationsNames = Annotations{
		Prefix:         "kube-ns-suspender",
		DesiredState:   "desiredState",
		ControllerName: "controllerName",
	}

	a.Clientset, err = createKubernetesClientset()
	if err != nil {
		return &a, err
	}
	return &a, nil
}

func createKubernetesClientset() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return nil, fmt.Errorf("cannot find kube config")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	config.WarningHandler = rest.NoWarnings{}

	// return the clientset
	return kubernetes.NewForConfig(config)
}

// UpdateNamespace updates a given namespace ns with a given state. If the state is not
// recognised or if the operation fails, it returns an error
func (a App) UpdateNamespace(ns, state string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		res, err := a.Clientset.CoreV1().Namespaces().Get(context.TODO(), ns, metav1.GetOptions{})
		if err != nil {
			return err
		}
		state := strings.Title(state)
		if state != RunningState && state != SuspendedState {
			return fmt.Errorf("state '%s' is not supported. Valid values are 'Running' and 'Suspended'", state)
		}
		// we set the annotation to the desired state
		res.Annotations[a.AnnotationsNames.Prefix+"/"+a.AnnotationsNames.DesiredState] = state
		_, err = a.Clientset.CoreV1().Namespaces().Update(context.TODO(), res, metav1.UpdateOptions{})
		return err
	})
}
