package pages

import (
	"embed"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/cypherfox/cloud-native-demo/pkg/k8s"
)

const failed_tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Welcome to BugSim</title>
	</head>
	<body>
	    <div>Dein niederträchtiger Angriff ist fehlgeschlagen!<p>Möchtest du noch einmal spielen?</div>
		<div><p>Zurück zur <a href="/">Startseite</a></p></div>
	</body>
</html>`

const pod_table_tpl = `
	<table>
	<tr><th><div>Name</div></th><th><div>Status</div></th><th><div>Alter</div></th></tr>
	{{range .Items}}
	    <tr>
		<td><div><a href="/api/delete/{{ .Name }}">{{ .Name }}</a></div></td>
		<td><div>{{ .State }}</div></td>
		<td><div>{{ .AgeString }}</div></td>
		</tr>
	{{else}}<div><strong>no pods</strong></div>{{end}}
	</table>
`

type PagesSetup struct {
	Namespace   string
	Deployment  string
	SuccessRate uint8
}

//go:embed root.tpl
var root_templ_file embed.FS
var root_templ *template.Template

var failed_templ *template.Template
var pod_table_templ *template.Template

// remember config params for easier access.
var setup PagesSetup

func Init(s PagesSetup) error {
	var err error

	// TODO: move this setup and the variable in deleteSinglePod.go to the K8s package.
	fmt.Println("Setting up Kubernetes Client")
	k8sClient, err = k8s.NewKubeClient()
	if err != nil {
		fmt.Printf("Initializing Kubernetes client failed: %s", err.Error())
		return err
	}

	root_templ_data, _ := root_templ_file.ReadFile("root.tpl")
	root_templ, err = template.New("rootPage").Parse(string(root_templ_data))
	if err != nil {
		fmt.Printf("Initializing root template failed: %s", err.Error())
		return err
	}
	failed_templ, err = template.New("failedPage").Parse(failed_tpl)
	if err != nil {
		fmt.Printf("Initializing failed attack template failed: %s", err.Error())
		return err
	}
	pod_table_templ, err = template.New("podDataTable").Parse(pod_table_tpl)
	if err != nil {
		fmt.Printf("Initializing pod data table template failed: %s", err.Error())
		return err
	}

	setup = s

	return nil
}

func respPrintf(w http.ResponseWriter, format string, a ...interface{}) error {
	_, err := io.WriteString(w, fmt.Sprintf(format, a...))
	if err != nil {
		fmt.Printf("cannot write to response: %s", err.Error())
		return err
	}
	return nil
}
