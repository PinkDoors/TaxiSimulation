package trip

type TripStatus string

const (
	DriverSearch TripStatus = "DRIVER_SEARCH"
	DriverFound  TripStatus = "DRIVER_FOUND" // accept
	OnPosition   TripStatus = "ON_POSITION"
	Started      TripStatus = "STARTED" // start
	Ended        TripStatus = "ENDED"
	Canceled     TripStatus = "CANCELED" // cancel
)
