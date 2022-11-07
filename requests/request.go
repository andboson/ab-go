package requests

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	METHOD_POST = "POST"
	METHOD_GET  = "GET"
)

type Request struct {
	Url      string
	Method   string
	Headers  []string
	PostData string
}

func (r *Request) Run(jobId string) *Response {
	var responseText string
	var contentLength int64
	var status int
	var reader io.ReadCloser
	var error error
	body := bytes.NewBuffer([]byte(r.PostData))
	request, error := http.NewRequest("POST", r.Url, body)
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	for _, header := range r.Headers {
		headersArray := strings.Split(header, ":")
		request.Header.Add(headersArray[0], headersArray[1])
	}
	HttpClient.Timeout = time.Duration(float64(DispatcherService.Timeout) * float64(time.Second))
	response, error := HttpClient.Do(request)

	if error == nil {
		defer response.Body.Close()
		status = response.StatusCode
		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			reader, error = gzip.NewReader(response.Body)
			defer reader.Close()
		default:
			reader = response.Body
		}

		var uncompressed []byte
		uncompressed, error = ioutil.ReadAll(reader)
		responseText = string(uncompressed)
		contentLength = int64(len(responseText))
	} else {
		log.Printf("\n Request error: %s \n request: %s", error, request)
	}

	if error != nil {
		status = 500
		responseText = error.Error()
	}

	return &Response{
		JobId:         jobId,
		Code:          status,
		RawResponse:   responseText,
		ContentLength: contentLength}
}
