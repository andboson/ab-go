package service
import "flag"

var Args *Flags

type Flags struct {
	Requests 	int
	Concurrency int
	Url			string
	PostData	string
	PostFile	string
	UrlFile		string
}

func readFlags(){
	flags := &Flags{}
	flags.PostData = *flag.Int("c", "", "a number")
	flag.Parse()
//	flag.Usage = func() {
//		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
//		flag.PrintDefaults()
//	}

	//url = os.Args[len(os.Args) - 1]
	//log.Printf("\n Sys args %s", url)

}
