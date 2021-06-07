package httpp

import "net/url"

type SampleInfo struct {
	URL        *url.URL
	Parameters url.Values
	Method     string
	IsRedirect bool
	FinalURL   *url.URL
}
