package requests

import (
	"time"
	"log"
	"os"
	"bufio"
	"net/http"
	"abgo/service"
	"math"
	"fmt"
)

var HttpClient *http.Client
var DispatcherService *Dispatcher

type Dispatcher struct {
	Args        *service.Flags
	Jobs        map[string]*Job
	Urls        []string
	PostData    []string
	ScannerPost *bufio.Scanner
	ScannerUrls *bufio.Scanner
	FilePtrPost *os.File
	FilePtrUrls *os.File
	Headers     []string
	Completed   []string
	Timeout		int
	Start 		time.Time
	Result		*Result
}

type Result struct {
	Duration 	string
	Requests	int
	Failed		int
	Rps			string
	Min 		string
	Max 		string
	Avg 		string

}

func CreateDispatcher() *Dispatcher{
	makeClient()
	DispatcherService = &Dispatcher{}
	DispatcherService.Result = &Result{}
	DispatcherService.Jobs = make(map[string]*Job)
	DispatcherService.Start = time.Now()
	DispatcherService.loadParams()

	return DispatcherService
}


func makeClient(){
	transport := &http.Transport{
		DisableCompression: false,
	}
	HttpClient = &http.Client{Transport: transport}
}

// run all processes
func (d *Dispatcher) Run() {
	defer d.FilePtrUrls.Close();
	defer d.FilePtrPost.Close();
	for i := 0; i < d.Args.Requests; i++ {
		var jobs = make(map[string]*Job)
		for j := 0; j < d.Args.Concurrency; j++ {
			if(i == d.Args.Requests){
				break
			}
			job := d.createJob()
			jobs[job.Id] = job
			i++
		}
		i--
		d.runBatch(jobs)
	}


	duration := time.Since(d.Start).Seconds()
	d.Result.Duration = fmt.Sprintf("%.3fms",duration * 1000)
	d.Result.Requests = len(d.Jobs)
	d.Result.Rps = fmt.Sprintf("%.0frps",math.Ceil(float64(d.Result.Requests)/duration))
	if(len(d.Jobs) < 1){
		return
	}

	var firstJob *Job
	for _, val := range( d.Jobs){
		firstJob = val
		break
	}
	min := firstJob.Duration
	max := firstJob.Duration
	var total float64
	for _, result := range d.Jobs{
		if( result.Duration > max){
			max = result.Duration
		}
		if( result.Duration < min){
			min = result.Duration
		}
		total = total + result.Duration
	}

	avg := total / float64(len(d.Jobs))
	d.Result.Avg = fmt.Sprintf("%.3fms", avg)
	d.Result.Max = fmt.Sprintf("%.3fms", max)
	d.Result.Min = fmt.Sprintf("%.3fms", min)

	return
}

//run async request in batch for specified amount of concurrency
//wait for all response
func (d *Dispatcher) runBatch(jobs map[string]*Job) {
	d.Completed = make([]string, 0)
	batchJobsCount := len(jobs)
	responseReciever := make(chan *Job, batchJobsCount)
	for _, job := range jobs {
		go job.Run(responseReciever)
	}

	func(){
		for {
			select {
			case response := <- responseReciever:
				d.Jobs[response.Id] = response
				if(response.Response.Code != 200){
					d.Result.Failed++
				}
				d.Completed = append(d.Completed, response.Id)
				if (len(d.Completed) == batchJobsCount) {
					return;
				}
			}
		}
	}()
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
	d.Headers = d.readHeaders()
}

//reads urls form file specified in -u param or from argument
func (d *Dispatcher) ReadUrl() string {
	if (d.Args.Url != "" && d.Args.UrlFile == "") {
		return d.Args.Url
	}

	if (d.FilePtrUrls == nil) {
		d.FilePtrUrls, _ = os.Open(d.Args.UrlFile)
		d.ScannerUrls = bufio.NewScanner(d.FilePtrUrls)
	}
	d.ScannerUrls.Split(bufio.ScanLines)
	d.ScannerUrls.Text()

	if (d.ScannerUrls.Scan()) {
		url := d.ScannerUrls.Text()
		if(url == ""){
			return d.ReadUrl()
		}

		return url
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
		text := d.ScannerPost.Text()
		if(text == ""){
			return d.ReadPostData()
		}

		return text
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
	if (d.Args.HeadersFile != "") {
		file, _ := os.Open(d.Args.HeadersFile)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			text := scanner.Text()
			d.Headers = append(d.Headers, text)
		}

	}
	var holder string
	for i, val := range d.Headers{
		if( holder == val){
			d.Headers = append(d.Headers[:i], d.Headers[i+1:]...)
		}
		holder = val
	}

	return d.Headers
}

//make job for request
func (d *Dispatcher) createJob() *Job {
	job := &Job{
		Id:service.RandStr(16, "number"),
		Completed:false,
		TimeStart:time.Now(),
		Request:d.makeRequest()}

	d.Jobs[job.Id] = job

	return job
}

//make request object from arguments
func (d *Dispatcher) makeRequest() *Request {
	requestObj := &Request{}
	method := METHOD_GET
	postData := d.ReadPostData()
	if (postData != "") {
		method = METHOD_POST;
	}
	url := d.ReadUrl()
	headers := d.Headers
	requestObj.Headers = headers
	requestObj.PostData = postData
	requestObj.Url = url
	requestObj.Method = method

	return requestObj
}
