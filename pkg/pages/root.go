package pages

import (
	"fmt"
	"net/http"
	"time"

	v1 "k8s.io/api/core/v1"
)

type podData struct {
	Name      string
	State     string
	AgeString string
}

func RootPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := respPrintf(w, "Hello, Bug from %s \n", r.RemoteAddr)
	if err != nil {
		return
	}

	podDataArr, err := getPodData()
	if err != nil {
		return
	}

	data := struct {
		Items []podData
	}{
		Items: *podDataArr,
	}

	err = root_templ.Execute(w, data)
	if err != nil {
		fmt.Printf("generating root page from template failed: %s\n", err.Error())
		return
	}

}

func getPodData() (*[]podData, error) {
	pods, err := k8sClient.GetPods(setup.Namespace, setup.Deployment)
	if err != nil {
		fmt.Printf("reading pods failed: %s\n", err.Error())
		return nil, err
	}

	podDataArr := []podData{}

	for _, pod := range pods.Items {
		podDataArr = append(podDataArr, podData{
			Name:      pod.GetName(),
			State:     statusMessage(pod),
			AgeString: time.Since(pod.GetCreationTimestamp().Time).String(),
		})
		fmt.Printf("%d: Setting %s to state %s, age %s",
			len(podDataArr),
			podDataArr[len(podDataArr)-1].Name,
			podDataArr[len(podDataArr)-1].State,
			podDataArr[len(podDataArr)-1].AgeString,
		)
	}
	return &podDataArr, nil
}

func statusMessage(pod v1.Pod) string {
	if pod.DeletionTimestamp != nil {
		return "Terminating"
	}
	return string(pod.Status.Phase)
}
