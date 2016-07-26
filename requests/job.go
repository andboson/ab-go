package requests

import (
	"time"
)

type Job struct {
	Id        string
	TimeStart time.Time
	Duration  float64
	Request   *Request
	Response  *Response
	Completed bool
}

func (j *Job) Run(resp chan *Job) {
	response := j.Request.Run(j.Id)
	j.Response = response
	timeDuration := time.Since(j.TimeStart)
	j.Duration = timeDuration.Seconds() * 1000
	j.Completed = true
	resp <- j
}
