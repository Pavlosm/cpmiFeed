package db

import "cpmiFeed/rawEventModels"

func NewEventDocument(event rawEventModels.Event) Event {
	return Event{
		Data:        event.Data,
		URL:         event.URL,
		Description: event.Description,
		Tags:        event.Tags,
		Timestamp:   event.Timestamp,
	}
}

func NewEventFromDocument(document Event) rawEventModels.Event {
	return rawEventModels.Event{
		Data:        document.Data,
		URL:         document.URL,
		Description: document.Description,
		Tags:        document.Tags,
		Timestamp:   document.Timestamp,
	}
}
