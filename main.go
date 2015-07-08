package main

import (
	"abgo/service"
	"log"
	"abgo/requests"
)

func  main(){
	service.Args.CheckUrl()
	requests.DispatcherService.Run();
	for _, result := range requests.DispatcherService.Jobs{
		log.Printf("\n code: %d \n result: %s", result.Response.Code, result.Duration)
	}

	log.Printf("\n results: %+v",requests.DispatcherService.Result)
}