package harkonnen

type Pecker struct {
	variables *map[string]interface{}

	controlChan chan<- string

	feedbackChan <-chan string

	script Script
}
