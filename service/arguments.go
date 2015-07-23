package service

import (
	"flag"
	"fmt"
	"os"
	"log"
	"regexp"
	"path/filepath"
	"strconv"
)

var Args *Flags

type Flags struct {
	Requests    int
	Concurrency int
	Timeout        int
	Url            string
	Header        string
	PostData    string
	PostFile    string
	UrlFile        string
	HeadersFile    string
	Port        string
	SlackUrl    string
	ApiName        string
	Testing        string
	Web            bool
	Ka             bool
}

func init() {
	curDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	match, _ := regexp.MatchString("_test", curDir)
	if (match) {
		return
	}
	ReadFlags()
}

func ReadFlags() {
	flags := &Flags{}
	//flags.Requests = *flag.Int("n", 1, "a number of requests")
	flag.IntVar(&flags.Requests, "n", 1, "a number of requests")
	flag.IntVar(&flags.Concurrency, "c", 1, "a number of concurrency requests")
	flag.StringVar(&flags.PostData, "d", "", "a string with post data")
	flag.StringVar(&flags.PostFile, "p", "", "a string, filename of file with postdata")
	flag.StringVar(&flags.UrlFile, "u", "", "a string, filename of file with urls")
	flag.StringVar(&flags.Header, "H", "", "a string, header")
	flag.StringVar(&flags.HeadersFile, "h", "", "a string, filename of file with headers")
	flag.StringVar(&flags.Port, "port", "9999", "a string, port of charts server")
	flag.StringVar(&flags.SlackUrl, "slack", "", "a string, Slack endpoint to send results")
	flag.StringVar(&flags.ApiName, "name", "app", "a string, tested api name")
	flag.IntVar(&flags.Timeout, "t", 3000, "a number, milliseconds request timeout")
	flag.StringVar(&flags.Testing, "test", "", "a string time duration of testing mode (repeat test)")
	//flag.BoolVar(&flags.Web,"web", false, "a flag, web mode (see localhost:9999 for results)")
	flag.BoolVar(&flags.Ka,"k", false, "a flag, Use HTTP KeepAlive feature")
	flags.Url = os.Args[len(os.Args) - 1]
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	Args = flags
}

func (f *Flags) CheckUrl() {
	if (f.UrlFile != "") {
		return
	}

	http, _ := regexp.MatchString("http", f.Url)
	if (len(f.Url) < 10 || !http) {
		log.Fatalf("Url incorrect!")
	}
}

func (f *Flags) GetDuration() float64 {
	var second int
	if (f.Testing == "0") {
		return 0
	}

	reg := regexp.MustCompile(`(\d+)(\w)`)
	match := reg.FindStringSubmatch(f.Testing)
	if (len(match) == 0) {
		intval, error := strconv.Atoi(f.Testing)
		if (error != nil) {
			log.Fatalf("Cannot get duration, %s", error)
		}

		return float64(intval)
	}

	duration, _ := strconv.Atoi(match[1])
	unit := match[2]

	switch unit {
	case "s":
		second = duration
	case "m":
		second = 60 * duration
	case "h":
		second = 3600 * duration
	}

	return float64(second)
}

