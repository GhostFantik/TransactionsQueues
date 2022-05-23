package repository

import (
	"TransactionsQueues/internal/models"
	"TransactionsQueues/pkg"
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (r Repository) CreateUser(name string) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(context.Background(), r.db, &user, "INSERT INTO users (name) VALUES ($1) RETURNING *", name)
	if err != nil {
		log.Err(err).Stack()
		return nil, pkg.NewError(http.StatusInternalServerError, err.Error())
	}
	return &user, nil
}

func (r Repository) GetUserByName(name string) (*models.User, error) {
	var user models.User
	err := pgxscan.Get(context.Background(), r.db, &user, "SELECT * FROM users WHERE name LIKE $1", name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.NewError(http.StatusNotFound, "not found")
		}
		log.Err(err).Stack()
		return nil, pkg.NewError(http.StatusInternalServerError, err.Error())
	}
	return &user, nil
}
