package requests

type Response struct {
	JobId       string
	Code        int
	RawResponse string
	ContentLength int64
}
