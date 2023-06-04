package k8s

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	NotImplementedYetErr = Error("function / feature not implemented yet")
	NotFoundErr          = Error("resource not found")
)

type K8sClient struct {
	Client *kubernetes.Clientset
}

func NewKubeClient() (*K8sClient, error) {
	var client K8sClient

	config, err := buildConfig()

	// creates the clientset
	client.Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func buildConfig() (*rest.Config, error) {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// assume out-of-cluster config when kubeconfig file exists
	_, err := os.Stat(*kubeconfig)

	var config *rest.Config
	if errors.Is(err, os.ErrNotExist) {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

	} else if err != nil {
		// other errors
		return nil, err

	} else {
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	return config, nil
}

func (kc *K8sClient) GetPods(namespace string, deplName string) (*coreV1.PodList, error) {

	listOpts := metaV1.ListOptions{}

	fmt.Printf("Selecting only pods for deployment %s\n", deplName)

	if deplName != "" {
		// Get deployment by name
		deployment, err := kc.getDeploymentByName(deplName, namespace)
		if err != nil {
			return nil, err
		}
		// read selector labels from deployment
		selector, err := metaV1.LabelSelectorAsSelector(deployment.Spec.Selector)
		if err != nil {
			return nil, err
		}
		listOpts.LabelSelector = selector.String()
	}

	// get pods in namespace and matching it against the label selector
	pods, err := kc.Client.CoreV1().Pods(namespace).List(context.TODO(), listOpts)
	if err != nil {
		fmt.Printf("reading pod info from cluster failed\n")
		return nil, err
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	return pods, nil
}

func (kc *K8sClient) getDeploymentByName(name, namespace string) (*appsV1.Deployment, error) {
	allDepls, err := kc.Client.AppsV1().Deployments(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, depl := range allDepls.Items {
		if depl.Name == name {
			return &depl, nil
		}
	}
	return nil, NotFoundErr
}

func (kc *K8sClient) DeletePod(name string, namespace string) error {
	return kc.Client.CoreV1().Pods(namespace).Delete(context.TODO(), name, metaV1.DeleteOptions{})
}
