package trip_outbound

type Event struct {
	ID              string `json:"id"`
	Source          string `json:"source"`
	Type            string `json:"type"`
	DataContentType string `json:"datacontenttype"`
	Time            string `json:"time"`
	Data            Data   `json:"data"`
}

type Data struct {
	TripID string        `json:"trip_id"`
	From   LatLngLiteral `json:"from"`
	To     LatLngLiteral `json:"to"`
	Price  Money         `json:"price"`
}

type LatLngLiteral struct {
	Lat float32 `bson:"lat"`
	Lng float32 `bson:"lng"`
}

type Money struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}
