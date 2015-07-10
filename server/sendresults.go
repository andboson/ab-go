package server
import (
	"ab-go/requests"
	"bytes"
	"net/http"
	"log"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"fmt"
	"ab-go/templates"
)

type Message struct {
	Text string `json:"text"`
}

func SendToSlack(dispatcher requests.Dispatcher) {
	message := Message{Text:fmt.Sprintf("\n Tested API: %s \n %s",
		dispatcher.Args.ApiName,
		templates.Formatter.FormatResult(dispatcher.Result))}
	bodyB, _ := json.Marshal(message)
	post := "payload=" + url.QueryEscape(string(bodyB))
	body := bytes.NewBuffer([]byte(post))
	request, err := http.NewRequest("POST", dispatcher.Args.SlackUrl, body)
	if ( err == nil) {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := requests.HttpClient.Do(request)
		if ( err != nil) {
			var reader io.ReadCloser
			reader = resp.Body
			var uncompressed []byte
			uncompressed, _ = ioutil.ReadAll(reader)
			responseText := string(uncompressed)
			log.Printf("\n errors send to slack,  \n error: %s \n response: %s", err, responseText)
		}

	} else {
		log.Printf("\n error send to slack", err)
	}
}
