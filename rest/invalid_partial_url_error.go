package rest

import "fmt"

type InvalidPartialUrlError struct {
	Url string
}

func (ipu InvalidPartialUrlError) Error() string {
	return fmt.Sprintf("'%s' is an invalid partial URL", ipu.Url)
}
