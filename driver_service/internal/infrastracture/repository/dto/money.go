package dto

type Money struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}
