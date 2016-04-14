# ab-go [![Build Status](https://travis-ci.org/andboson/ab-go.svg?branch=master)](https://travis-ci.org/andboson/ab-go)

Inspired by `apache ab testing tool`. Created with `golang`.
 
 
Features:

- File with post data strings
- File with urls strings
- Headers file
- Real-time stats
- Sends result to Slack


run ./ab-go for usage

## run options:

##### -n
Number of requests sended to server.

##### -c
Number of concurrency requests in one batch.

##### -d
String with post data

##### -p
String with filename contains post data strings

##### -u
String with filename contains urls strings

##### -H
String with header

##### -h
String with filename contains headers strings

##### -t
Number of milliseconds request timeout

##### -test
String with time duration. Enable cycled tests with sending results to graphics.<br>
Sample values: "5m", "60s", "24h". Value "0" starts endless testing mode.<br>
Run with this flag and open `localhost:9999` (default port) in your browser<br>
<img src="tests/screenshot.png" />
Data is updated every second.

##### -port
Embedded web-server port. Used with `-test` flag

##### -slack
String with Slack endpoint url (incoming WebHoock) for sending results.

##### -name
String of tested api name, who will be presented in results message in Slack. Used with `-slack` flag

##### -k
Use HTTP KeepAlive feature

-----------------
#### build:
mkdir $GOPATH/src/github.com/andboson/ab-go && cd $GOPATH/src/github.com/andboson/ab-go

git clone git@github.com:andboson/ab-go.git .

make
