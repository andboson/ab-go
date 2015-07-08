package requests

import (
	"math/rand"
	"time"
	"log"
	"os"
	"bufio"
	"net/http"
	"abgo/service"
)

var HttpClient *http.Client
var DispatcherService *Dispatcher

type Dispatcher struct {
	Args        *service.Flags
	Jobs        map[int32]*Job
	Urls        []string
	PostData    []string
	ScannerPost *bufio.Scanner
	ScannerUrls *bufio.Scanner
	FilePtrPost *os.File
	FilePtrUrls *os.File
	Headers     []string
	Completed   []int32
	Timeout		int
}

func init() {
	makeClient()
	DispatcherService = &Dispatcher{}
	DispatcherService.Jobs = make(map[int32]*Job)
	DispatcherService.loadParams()
}

func makeClient(){
	transport := &http.Transport{
		DisableCompression: false,
	}
	HttpClient = &http.Client{Transport: transport}
}

// run all processes
func (d *Dispatcher) Run() {
	var jobs = make(map[int32]*Job)
	defer d.FilePtrUrls.Close();
	defer d.FilePtrPost.Close();
	for i := 0; i < d.Args.Requests; i++ {
		for j := 0; j < d.Args.Concurrency; j++ {
			jobs = d.makeJob()
		}
		d.runBatch(jobs)
	}
}

//run async request in batch for specified amount of concurrency
//wait for all response
func (d *Dispatcher) runBatch(jobs map[int32]*Job) {
	d.Completed = make([]int32, 0)
	batchJobsCount := len(jobs)
	responseReciever := make(chan *Job, batchJobsCount)
	completedReciever := make(chan int32, batchJobsCount)
	for _, job := range jobs {
		go job.Run(responseReciever, completedReciever)
	}

	for {
		select {
		case response := <- responseReciever:
			d.Jobs[response.Id] = response

		case completed := <-completedReciever:
			d.Completed = append(d.Completed, completed)
			if (len(d.Completed) == batchJobsCount) {
				break
			}
		}
	}

}

//load all params
func (d *Dispatcher) loadParams() {
	d.Args = service.Args
	if (d.Args.UrlFile != "") {
		return
	} else if (d.Args.Url != "") {
		d.Urls = append(d.Urls, d.Args.Url)
	} else {
		log.Fatalf("You must specify at once one url! ", d.Args.Url)
	}
	d.Timeout = d.Args.Timeout
	d.readHeaders()
}

//reads urls form file specified in -u param or from argument
func (d *Dispatcher) ReadUrl() string {
	if (d.Args.Url != "") {
		return d.Args.Url
	}

	if (d.FilePtrUrls == nil) {
		d.FilePtrUrls, _ = os.Open(d.Args.UrlFile)
		d.ScannerUrls = bufio.NewScanner(d.FilePtrUrls)
	}
	d.ScannerUrls.Split(bufio.ScanLines)
	d.ScannerUrls.Text()

	if (d.ScannerUrls.Scan()) {
		return d.ScannerUrls.Text()
	} else {
		d.FilePtrUrls.Close();
		d.FilePtrUrls = nil
		return d.ReadUrl()
	}
}

//reads postdata form file specified in -p param or string from -d flag
func (d *Dispatcher) ReadPostData() string {
	if (d.Args.PostData != "") {
		return d.Args.PostData
	}

	if (d.Args.PostFile == "") {
		return "";
	}

	if (d.FilePtrPost == nil) {
		d.FilePtrPost, _ = os.Open(d.Args.PostFile)
		d.ScannerPost = bufio.NewScanner(d.FilePtrPost)
	}
	d.ScannerPost.Split(bufio.ScanLines)
	d.ScannerPost.Text()

	if (d.ScannerPost.Scan()) {
		return d.ScannerPost.Text()
	} else {
		d.FilePtrPost.Close();
		d.FilePtrPost = nil
		return d.ReadPostData()
	}
}

//reads headers form file specified in -h param or string from -H param
func (d *Dispatcher) readHeaders() []string{
	if (d.Args.Header != "") {
		d.Headers = append(d.Headers, d.Args.Header)
	}
	if (d.Args.HeadersFile == "") {
		return d.Headers
	}
	file, _ := os.Open(d.Args.UrlFile)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		d.Headers = append(d.Headers, scanner.Text())
	}

	return d.Headers
}

//make job for request
func (d *Dispatcher) makeJob() map[int32]*Job {
	var jobs = make(map[int32]*Job)
	job := &Job{
		Id:rand.Int31n(10),
		Completed:false,
		TimeStart:time.Now(),
		Request:d.makeRequest()}

	d.Jobs[job.Id] = job
	jobs[job.Id] = job

	return jobs
}

//make request object from arguments
func (d *Dispatcher) makeRequest() *Request {
	requestObj := &Request{}
	method := METHOD_GET
	postData := d.ReadPostData()
	if (postData != "") {
		method = METHOD_POST;
	}
	requestObj.Headers = d.readHeaders()
	requestObj.PostData = postData
	requestObj.Url = d.ReadUrl()
	requestObj.Method = method

	return requestObj
}
