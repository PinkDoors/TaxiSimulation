package users

type Balance struct {
	Amount float64
}

type User struct {
	ID        int64    `json:"id"`
	FirstName string   `json:"first-name"`
	LastName  string   `json:"last-name"`
	Balance   *Balance `json:"balance"`
}

type Users []User
