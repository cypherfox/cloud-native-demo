package pages

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cypherfox/cloud-native-demo/pkg/version"
	"github.com/hako/durafmt"
	v1 "k8s.io/api/core/v1"
)

type podData struct {
	Name      string
	State     string
	AgeString string
}

func RootPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	/*
		err := respPrintf(w, "Hello, Bug from %s \n", r.RemoteAddr)
		if err != nil {
			return
		}
	*/

	podDataArr, err := getPodData()
	if err != nil {
		return
	}

	data := struct {
		Items       []podData
		Version     string
		SuccessRate uint8
	}{
		Items:       *podDataArr,
		Version:     version.BuildVersion,
		SuccessRate: setup.SuccessRate,
	}

	err = root_templ.Execute(w, data)
	if err != nil {
		fmt.Printf("generating root page from template failed: %s\n", err.Error())
		return
	}

}

func StyleSheet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")

	w.Write(style_css)
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
			AgeString: prettyPrintTime(time.Since(pod.GetCreationTimestamp().Time)),
		})

		/*
			fmt.Printf("%d: Setting %s to state %s, age %s\n",
				len(podDataArr),
				podDataArr[len(podDataArr)-1].Name,
				podDataArr[len(podDataArr)-1].State,
				podDataArr[len(podDataArr)-1].AgeString,
			)
		*/
	}
	return &podDataArr, nil
}

func statusMessage(pod v1.Pod) string {
	if pod.DeletionTimestamp != nil {
		return "Terminating"
	}
	return string(pod.Status.Phase)
}

func prettyPrintTime(duration time.Duration) string {
	durStr := durafmt.Parse(duration).LimitFirstN(3).String()

	return durStr
}
