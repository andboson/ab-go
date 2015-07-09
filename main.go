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
	run(false)
	if(service.Args.Tesing){
		for{
			run(true)
		}
	}
}

func run(clearScreen bool){
	dispatcher := requests.CreateDispatcher()
	dispatcher.Run()
	if(clearScreen){
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	}
	fmt.Printf("\n %s", templates.Formatter.FormatResult(requests.DispatcherService.Result))
}