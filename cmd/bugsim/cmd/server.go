/*
Copyright © 2022 Lutz Behnke <lutz.behnke@gmx.de>
This file is part of cloud-native demo
*/
package cmd

import (
	"fmt"

	"log"
	"net/http"
	"os"

	mux "github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cypherfox/cloud-native-demo/pkg/pages"
)

var Port uint16

var Namespace string
var Deployment string
var SuccessRate uint8

var router *mux.Router

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "run the web listener",
	Long:  `run the web listener`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("server called")

		err := doServer()
		if err != nil {
			return err
		}
		fmt.Println("server ended")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Uint16VarP(&Port, "port", "p", 80, "Port on which to listen for web requests")
	serverCmd.Flags().StringVarP(&Namespace, "namespace", "n", "default", "Namespace in which to look for pods.")
	serverCmd.Flags().StringVarP(&Deployment, "deployment", "d", "web", "Deployment from which to delete pods.")
	serverCmd.Flags().Uint8VarP(&SuccessRate, "success-rate", "r", 15, "Success rate of the simulated bug attack on the pods in percent [1-100] (default are 15%")
}

func doServer() error {
	var err error

	if SuccessRate < 1 || SuccessRate > 100 {
		fmt.Printf("Invalid success rate percentage: %d. Must be between 1 and 100", SuccessRate)
		return os.ErrInvalid
	}

	err = pages.Init(pages.PagesSetup{
		Namespace:   Namespace,
		Deployment:  Deployment,
		SuccessRate: SuccessRate,
	})
	if err != nil {
		return err
	}

	router = mux.NewRouter()

	fmt.Println("Setting up pages")
	router.HandleFunc("/", pages.RootPage)
	router.HandleFunc("/api/delete/{id}", pages.DeleteSinglePod)
	// router.HandleFunc("/ws/pod_info", pages.PodInfoReloaderWS)
	http.Handle("/", router)

	fmt.Printf("Starting to serve on port %d\n", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))

	return nil
}
