package models

type LatLngLiteral struct {
	Lat float32 `bson:"lat"`
	Lng float32 `bson:"lng"`
}
