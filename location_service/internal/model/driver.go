package model

type Driver struct {
	DriverId string  `json:"id"`
	Lat      float32 `json:"lat"`
	Lng      float32 `json:"lng"`
}
