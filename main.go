package main

import (
	"fmt"
	"github.com/cosmic-chichu/ab-go/requests"
	"github.com/cosmic-chichu/ab-go/server"
	"github.com/cosmic-chichu/ab-go/service"
	"github.com/cosmic-chichu/ab-go/templates"
	"os"
	"os/exec"
	"time"
)

func main() {
	//go server.Init()
	service.Args.CheckUrl()
	run(false, false)
	if service.Args.Testing != "" {
		run(true, true)
	}
}

func run(clearScreen bool, testing bool) {
	if clearScreen {
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	}

	if service.Args.Testing != "" {
		var timeout <-chan time.Time
		exit := false
		duration := service.Args.GetDuration()
		if duration != 0 {
			timeout = time.After(time.Duration(duration * float64(time.Second)))
		}
		ticker := time.NewTicker(time.Second)
		dispatcher := requests.CreateDispatcher()
		for {
			select {
			case <-timeout:
				exit = true
			case <-ticker.C:
				if exit {
					return
				}
				dispatcher.Run()
				c := exec.Command("clear")
				c.Stdout = os.Stdout
				c.Run()
				fmt.Printf("#AB-GO testing tool. \n\n Testing %s...\n Open http://localhost:%s for results \n Current results: \n %s",
					service.Args.ApiName,
					service.Args.Port,
					templates.Formatter.FormatResult(dispatcher.Result))

				fmt.Printf("\n\n\n Last response: \n\n %s", dispatcher.Result.LastResult)

				server.Send <- dispatcher.Result
			}
		}
	} else {
		dispatcher := requests.CreateDispatcher()
		dispatcher.Run()
		// fmt.Printf("#AB-GO testing tool. \n\n Testing %s.  \n Results: \n %s",
		// 	service.Args.ApiName,
		// 	templates.Formatter.FormatResult(dispatcher.Result))
		fmt.Println("Requests, Failed, Duration, RPS, Min, Max, Avg")
		fmt.Printf("%d, %d, %s, %s, %s, %s, %s", dispatcher.Result.Requests, dispatcher.Result.Failed, dispatcher.Result.Duration, dispatcher.Result.Rps, dispatcher.Result.Min,
		dispatcher.Result.Max, dispatcher.Result.Avg)

		if dispatcher.Args.SlackUrl != "" {
			server.SendToSlack(*dispatcher)
		}
	}
}
