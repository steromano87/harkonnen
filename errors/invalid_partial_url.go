package errors

import "fmt"

type InvalidPartialUrl struct {
	Url string
}

func (ipu InvalidPartialUrl) Error() string {
	return fmt.Sprintf("'%s' is an invalid partial URL", ipu.Url)
}
