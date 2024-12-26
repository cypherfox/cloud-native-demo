package pages

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/cypherfox/cloud-native-demo/pkg/k8s"

	"github.com/gorilla/mux"
)

var k8sClient *k8s.K8sClient

func DeleteSinglePod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podName := vars["id"]

	// Cast a vote
	probability := rand.Float32()
	if probability > (float32(setup.SuccessRate) / 100.0) {
		data := struct{}{}

		err := failed_templ.Execute(w, data)
		if err != nil {
			fmt.Printf("generating root page from template failed: %s\n", err.Error())
		}

		return
	}

	// vote was successful -> kill the pod
	err := k8sClient.DeletePod(podName, setup.Namespace)
	if err != nil {
		fmt.Printf("deleting pod %s failed: %s\n", podName, err.Error())
		return
	}

	// TODO: change this into a redirect, in order to clean up the URL.
	RootPage(w, r)
}
