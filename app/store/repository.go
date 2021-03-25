package store

import "github.com/iavealokin/cashflow/app/model"

type UserRepository interface {
	Create(*model.Operation) error
	Drop(*model.User) error
	Update(*model.User) error
	GetOperations(int) ([]model.Operation, error)
	GetUserData(int) (*model.UserData, error)
	UserLogin(string, string) (*model.User, error)
}
