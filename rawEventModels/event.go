package rawEventModels

type Event struct {
	Data        string   `json:"data"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
