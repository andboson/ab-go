package service
import (
	"flag"
	"fmt"
	"os"
	"log"
	"regexp"
	"path/filepath"
)

var Args *Flags

type Flags struct {
	Requests 	int
	Concurrency int
	Timeout		int
	Url			string
	Header		string
	PostData	string
	PostFile	string
	UrlFile		string
	HeadersFile	string
}

func init(){
	curDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	match, _ := regexp.MatchString("_test",curDir)
	if(match){
		return
	}
	readFlags()
}

func readFlags(){
	flags := &Flags{}
	flags.Requests = *flag.Int("n", 1, "a number of requests")
	flags.Concurrency = *flag.Int("c", 1, "a number of concurrency requests")
	flags.PostData = *flag.String("d", "", "a string with post data")
	flags.PostFile = *flag.String("p", "", "a string, filename of file with postdata")
	flags.UrlFile = *flag.String("u", "", "a string, filename of file with urls")
	flags.Header = *flag.String("H", "", "a string, header")
	flags.HeadersFile = *flag.String("h", "", "a string, filename of file with headers")
	flags.Timeout = *flag.Int("t", 3000, "a number, milliseconds request timeout")
	flags.Url = os.Args[len(os.Args) - 1]
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	Args = flags
}

func (f *Flags) CheckUrl(){
	if(f.UrlFile != ""){
		return
	}

	http, _ := regexp.MatchString("http",f.Url)
	if(len(f.Url) < 10 || !http){
		log.Fatalf("Url incorrect!")
	}
}
