package requests

import (
	"time"
)

type Job struct {
	Id			int32
	TimeStart	time.Time
	TimeEnd		time.Time
	Duration	time.Duration
	Request		*Request
	Response	*Response
	Completed	bool
}


func (j *Job) Run(resp chan *Job, completed chan int32){
	response := j.Request.Run(j.Id)
	j.Response = response
	resp<-j
	completed<-j.Id
}