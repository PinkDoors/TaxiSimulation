package trip_inbound

// Event представляет собой структуру для хранения информации об эвенте из Kafka.
type Event struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	//Type            EventType `json:"type"`
	Type            string `json:"type"`
	DataContentType string `json:"datacontenttype"`
	Time            string `json:"time"`
	Data            Data   `json:"data"`
}

// EventType представляет собой тип события.
//type EventType string
//
//// Enum для типа события.
//const (
//	TripEventAccepted EventType = "trip.event.accepted"
//)

// Data представляет собой тело события.
type Data struct {
	TripID string `json:"trip_id"`
}
