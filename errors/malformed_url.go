package errors

import "fmt"

type MalformedUrl struct {
	Url string
}

func (mu MalformedUrl) Error() string {
	return fmt.Sprintf("'%s' is a malformed URL", mu.Url)
}
