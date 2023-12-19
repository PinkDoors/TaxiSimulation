package users

import (
	"Task3/internal/domain/models/user"
)

func ToModelsUser(userToParse User) *user.User {
	return &user.User{
		Id:        userToParse.ID,
		FirstName: userToParse.FirstName,
		LastName:  userToParse.LastName,
		Balance: &user.Balance{
			Amount: userToParse.Balance.Amount,
		},
	}
}
