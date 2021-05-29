package cmd

import (
	"accurics-test/config"
	"accurics-test/router"
	"flag"
	"fmt"
	"net/http"
)

const (
	defaultConfigFile = "config.json"
)

//RootCmd is default function which should be called when execution begins
func RootCmd() error {
	fmt.Println("service started!")
	configFile := flag.String("config", defaultConfigFile, "path to json configuration file")
	flag.Parse()
	err := config.LoadApplication(*configFile)
	if err != nil {
		fmt.Println("error in loading App = ", err)
	}
	r := router.NewRouter()
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error occured while listening :", err)
	}
	return nil
}
