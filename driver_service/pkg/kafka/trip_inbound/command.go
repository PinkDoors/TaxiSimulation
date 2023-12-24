package trip_inbound

import (
	"github.com/google/uuid"
	"time"
)

type Command struct {
	ID              uuid.UUID `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	DataContentType string    `json:"datacontenttype"`
	Time            time.Time `json:"time"`
	Data            Data      `json:"data"`
}

type Data struct {
	TripID   uuid.UUID `json:"trip_id"`
	DriverID uuid.UUID `json:"driver_id"`
}
