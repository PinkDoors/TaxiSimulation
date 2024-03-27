package trip

type Status string

const (
	DriverSearch Status = "DRIVER_SEARCH"
	DriverFound  Status = "DRIVER_FOUND" // accept
	OnPosition   Status = "ON_POSITION"
	Started      Status = "STARTED" // start
	Ended        Status = "ENDED"
	Canceled     Status = "CANCELED" // cancel
)
