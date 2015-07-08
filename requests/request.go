package requests
import (
	"net/http"
	"bytes"
	"time"
	"compress/gzip"
	"io/ioutil"
	"io"
)

const (
	METHOD_POST="POST"
	METHOD_GET="GET"
)

type Request struct {
	Url			string
	Method		string
	Headers		[]string
	PostData	string
}

func (r *Request) Run(jobId int32) *Response{
	var responseText string
	var status int
	var reader io.ReadCloser
	body := bytes.NewBuffer([]byte(r.PostData))
	request, error := http.NewRequest("POST", r.Url, body)
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	HttpClient.Timeout = time.Duration(float64(DispatcherService.Timeout) * float64(time.Second))
	response, error := HttpClient.Do(request)
	defer response.Body.Close()
	status = response.StatusCode

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, error = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}

	uncompressed, error := ioutil.ReadAll(reader)
	responseText = string(uncompressed)

	if(error != nil){
		status = 500
		responseText = error.Error()
	}

	return &Response{
		JobId:jobId,
		Code:status,
		RawResponse:responseText}
}
