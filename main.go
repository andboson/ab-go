package main

import (
	"fmt"
	"github.com/andboson/ab-go/service"
	"github.com/andboson/ab-go/requests"
	"github.com/andboson/ab-go/templates"
	"os/exec"
	"os"
	"github.com/andboson/ab-go/server"
	"time"
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
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				dispatcher := requests.CreateDispatcher()
				dispatcher.Run()
				c := exec.Command("clear")
				c.Stdout = os.Stdout
				c.Run()
				fmt.Printf("\n %s", templates.Formatter.FormatResult(dispatcher.Result))
				server.Send<-dispatcher.Result
			}
		}
	} else {
		dispatcher := requests.CreateDispatcher()
		dispatcher.Run()
		fmt.Printf("\n %s", templates.Formatter.FormatResult(dispatcher.Result))
		if(dispatcher.Args.SlackUrl != ""){
			server.SendToSlack(*dispatcher)
		}
	}
}
