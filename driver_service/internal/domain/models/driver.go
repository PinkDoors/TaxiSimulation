package models

import "github.com/google/uuid"

type Driver struct {
	Id       uuid.UUID
	Location LatLngLiteral
}
