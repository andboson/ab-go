package main
import (
	"log"
	//"os"
	"flag"
	"os"
	"fmt"
)

func main(){
	wordPtr := flag.String("word", "", "a string")
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	log.Printf("\n Sys args %s", *wordPtr)
}