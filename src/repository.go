package gotodo

import (
	"database/sql"
	"fmt"
)

type IRepository[T any] interface {
	Create(*T) (*T, error)
	Update(*T) (*T, error)
	Delete(int) (bool, error)
	All() ([]*T, error)
	Get(int) (*T, error)
}

type BaseRepository struct {
	*sql.DB
}

func CreateBaseRepository(config *Config) (*BaseRepository, error) {
	connStr := fmt.Sprintf("%s:%s@/%s", config.MySql.UserName, config.MySql.Password, config.MySql.DatabaseName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	repository := BaseRepository{
		DB: db,
	}

	return &repository, nil
}
