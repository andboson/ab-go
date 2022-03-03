package templates

import (
	"fmt"
	"github.com/cosmic-chichu/ab-go/requests"
)

var Formatter *Format
var ResultString = `
	Requests:		 %d
	Failed requests:	 %d
	Duration:		 %s
	Rps:			 %s
	Min:			 %s
	Max:			 %s
	Avg:			 %s`

type Format struct {
	*string
}

func init() {
	Formatter = &Format{}
}

func (f *Format) FormatResult(result *requests.Result) string {
	string := fmt.Sprintf(ResultString, result.Requests, result.Failed, result.Duration, result.Rps, result.Min,
		result.Max, result.Avg)

	return string
}
