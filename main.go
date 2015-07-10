package main

import (
	"fmt"
	"abgo/service"
	"abgo/requests"
	"abgo/templates"
	"os/exec"
	"os"
	"abgo/server"
)

func  main(){
	go server.Init()
	service.Args.CheckUrl()
	run(false, false)
	if(service.Args.Tesing){
		run(true, true)
	}
}

func run(clearScreen bool, testing bool){
	if (clearScreen) {
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	}
	if(testing) {
		for {
			dispatcher := requests.CreateDispatcher()
			dispatcher.Run()
			c := exec.Command("clear")
			c.Stdout = os.Stdout
			c.Run()
			fmt.Printf("\n %s", templates.Formatter.FormatResult(requests.DispatcherService.Result))
			server.Send<-requests.DispatcherService.Result
		}
	} else {
		dispatcher := requests.CreateDispatcher()
		dispatcher.Run()
		fmt.Printf("\n %s", templates.Formatter.FormatResult(requests.DispatcherService.Result))
	}
}