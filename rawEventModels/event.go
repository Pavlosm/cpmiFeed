package rawEventModels

import "time"

type Event struct {
	Data        string    `json:"data"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Timestamp   time.Time `json:"timestamp"`
}
