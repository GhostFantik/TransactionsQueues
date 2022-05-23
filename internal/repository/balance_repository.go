package repository

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (r Repository) PutBalance(userId uuid.UUID, money int) error {
	err := r.PushTask(userId, PutMoneyOperation, money)
	if err != nil {
		log.Err(err).Stack().Msg("ошибка!")
		return err
	}
	return nil
}

func (r Repository) PopBalance(userId uuid.UUID, money int) error {
	err := r.PushTask(userId, PopMoneyOperation, money)
	if err != nil {
		log.Err(err).Stack().Msg("ошибка!")
		return err
	}
	return nil
}
