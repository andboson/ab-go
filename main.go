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

func main() {
	go server.Init()
	service.Args.CheckUrl()
	//run(false, false)
	//	if(service.Args.Testing){
	run(true, true)
	//	}
}

func run(clearScreen bool, testing bool) {
	if (clearScreen) {
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	}

	if (service.Args.Testing != "") {
		var timeout <- chan time.Time
		exit := false
		duration := service.Args.GetDuration()
		if (duration != 0) {
			timeout = time.After(time.Duration(duration * float64(time.Second)))
		}
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-timeout:
				exit = true
			case <-ticker.C:
				if (exit) {
					return
				}
				dispatcher := requests.CreateDispatcher()
				dispatcher.Run()
				c := exec.Command("clear")
				c.Stdout = os.Stdout
				c.Run()
				fmt.Printf("#AB-GO testing tool. \n\n Testing %s...\n Open http://localhost:%s for results \n Current results \n %s",
					service.Args.ApiName,
					service.Args.Port,
					templates.Formatter.FormatResult(dispatcher.Result))

				server.Send <- dispatcher.Result
			}
		}
	} else {
		dispatcher := requests.CreateDispatcher()
		dispatcher.Run()
		fmt.Printf("\n %s", templates.Formatter.FormatResult(dispatcher.Result))
		if (dispatcher.Args.SlackUrl != "") {
			server.SendToSlack(*dispatcher)
		}
	}
}
