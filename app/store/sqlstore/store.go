package sqlstore

import (
	"database/sql"

	"github.com/iavealokin/cashflow/app/store"
	_ "github.com/lib/pq" // ....
)

//Store ..
type Store struct {
	db *sql.DB
	UserRepository *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository{
	if s.UserRepository!= nil{
		return s.UserRepository
	}

	s.UserRepository=&UserRepository{
		store: s,
	}
		return s.UserRepository
}