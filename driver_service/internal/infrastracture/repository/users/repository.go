package users

import (
	errors2 "Task3/internal/domain/errors"
	"Task3/internal/domain/models/user"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type UserJsonFileRepository struct {
	logger   *zap.Logger
	FilePath string
}

func NewUserJsonFileRepository(
	logger *zap.Logger,
	filePath string,
) *UserJsonFileRepository {
	return &UserJsonFileRepository{
		logger:   logger,
		FilePath: filePath,
	}
}

func (r UserJsonFileRepository) CreateUser(ctx context.Context, firstName string, lastName string) error {
	file, err := os.Open(r.FilePath)
	if err != nil {
		r.logger.Error("Error during attempt to open users file:", zap.Error(err))
		return err
	}
	defer file.Close()

	var users Users
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		r.logger.Error("Error when trying to read users file:", zap.Error(err))
		return err
	}

	var maxId int64 = 0
	for _, currentUser := range users {
		if currentUser.ID > maxId {
			maxId = currentUser.ID
		}
	}

	newId := maxId + 1

	users = append(users, User{
		newId,
		firstName,
		lastName,
		&Balance{
			0,
		},
	})

	file, err = os.Create(r.FilePath)
	if err != nil {
		r.logger.Error("Error during attempt to open users file:", zap.Error(err))
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(users); err != nil {
		r.logger.Error("Error encoding JSON:", zap.Error(err))
		return err
	}

	return err
}

func (r UserJsonFileRepository) GetUserBalance(ctx context.Context, userId int64) (*user.Balance, error) {
	file, err := os.Open(r.FilePath)
	if err != nil {
		r.logger.Error("Error during attempt to open users file:", zap.Error(err))
		return nil, err
	}
	defer file.Close()

	var users Users
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		r.logger.Error("Error when trying to read users file:", zap.Error(err))
		return nil, err
	}

	for _, currentUser := range users {
		if currentUser.ID == userId {
			return ToModelsUser(currentUser).Balance, nil
		}
	}

	return nil, errors2.NotFoundError{Message: "User not found for id: " + strconv.FormatInt(userId, 10)}
}

func (r UserJsonFileRepository) IncreaseBalance(ctx context.Context, userId int64, amount float64) error {
	file, err := os.Open(r.FilePath)
	if err != nil {
		r.logger.Error("Error during attempt to open users file:", zap.Error(err))
		return err
	}
	defer file.Close()

	var users Users
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		r.logger.Error("Error when trying to read users file:", zap.Error(err))
		return err
	}

	var user *User = nil
	for i, currentUser := range users {
		if currentUser.ID == userId {
			user = &users[i]
			break
		}
	}

	if user == nil {
		return errors2.NotFoundError{Message: "User not found for id: " + strconv.FormatInt(userId, 10)}
	}

	user.Balance.Amount += amount

	file, err = os.Create(r.FilePath)
	if err != nil {
		r.logger.Error("Error during attempt to open users file:", zap.Error(err))
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(users); err != nil {
		r.logger.Error("Error encoding JSON:", zap.Error(err))
		return err
	}

	return nil
}

func (r UserJsonFileRepository) TransferBetweenBalances(
	ctx context.Context,
	senderUserId int64,
	receiverUserId int64,
	amount float64) error {
	file, err := os.Open(r.FilePath)
	if err != nil {
		r.logger.Error("Error during attempt to open users file:", zap.Error(err))
		return err
	}
	defer file.Close()

	var users Users
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		r.logger.Error("Error when trying to read users file:", zap.Error(err))
		return err
	}

	var receiverUser *User = nil
	var senderUser *User = nil

	for i, currentUser := range users {
		if currentUser.ID == receiverUserId {
			receiverUser = &users[i]
		}
		if currentUser.ID == senderUserId {
			senderUser = &users[i]
		}
	}

	if receiverUser == nil {
		return errors2.NotFoundError{Message: "User not found for id: " + strconv.FormatInt(receiverUserId, 10)}
	}
	if senderUser == nil {
		return errors2.NotFoundError{Message: "User not found for id: " + strconv.FormatInt(senderUserId, 10)}
	}

	if senderUser.Balance.Amount-amount < 0 {
		return errors2.NegativeBalanceError{
			Message: "User with id: " + strconv.FormatInt(receiverUserId, 10) +
				" will have negative balance after transaction"}
	}

	senderUser.Balance.Amount -= amount
	receiverUser.Balance.Amount += amount

	file, err = os.Create(r.FilePath)
	if err != nil {
		r.logger.Error("Error during attempt to open users file:", zap.Error(err))
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(users); err != nil {
		r.logger.Error("Error encoding JSON:", zap.Error(err))
		return err
	}

	return err
}
