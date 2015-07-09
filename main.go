package main

import (
	"fmt"
	"abgo/service"
	"abgo/requests"
	"abgo/templates"
)

func  main(){
	service.Args.CheckUrl()
	requests.DispatcherService.Run();
//	for _, result := range requests.DispatcherService.Jobs{
//		log.Printf("\n code: %d \n result: %s", result.Response.Code, result.Duration)
//	}
	fmt.Println(templates.Formatter.FormatResult(requests.DispatcherService.Result))
}