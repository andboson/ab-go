package main

import (
	"abgo/service"
	"log"
	"abgo/requests"
)

func  main(){
	service.Args.CheckUrl()
	requests.DispatcherService.Run();
	log.Printf("\n =====", requests.DispatcherService.Jobs)
}