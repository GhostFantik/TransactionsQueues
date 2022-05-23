package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

const (
	PutMoneyOperation = "+"
	PopMoneyOperation = "-"
)

var clientsChannels map[uuid.UUID]chan int = map[uuid.UUID]chan int{}

func (r Repository) clientsQueue(userId uuid.UUID) {
	stopTimer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-stopTimer.C:
			_, err := r.rdb.LRem(context.Background(), "active_clients", 0, userId.String()).Result()
			if err != nil {
				log.Err(err).Stack().Msg("ошибка Redis")
			}
			delete(clientsChannels, userId)
			log.Info().Stack().Msg(userId.String())
			return
		default:
			// for debug
			time.Sleep(5 * time.Second)

			operation, err := r.rdb.LPop(context.Background(), userId.String()).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				log.Err(err).Stack().Msg("ошибка получения данных из Redis")
			}
			if len(operation) < 0 {
				continue
			}
			money, err := strconv.Atoi(operation[1:])
			if err != nil {
				log.Error().Msg("ошибка данных в Redis")
			}
			stopTimer.Reset(10 * time.Second)
			if operation[0:1] == PutMoneyOperation {
				_, err := r.db.Exec(context.Background(), "UPDATE users "+
					"SET balance = ((SELECT balance FROM users WHERE id = $1) + $2) "+
					"WHERE id = $1",
					userId, money)
				if err != nil {
					log.Err(err).Stack().Msg("ошибка обновления баланса")
				}
			} else if operation[0:1] == PopMoneyOperation {
				var balance int
				err := pgxscan.Get(context.Background(), r.db, &balance, "SELECT balance FROM users WHERE id = $1", userId)
				if err != nil {
					log.Err(err).Stack().Msg("ошибка получения данных из Redis")
				}
				if balance-money >= 0 {
					_, err := r.db.Exec(context.Background(), "UPDATE users "+
						"SET balance = ((SELECT balance FROM users WHERE id = $1) - $2) "+
						"WHERE id = $1",
						userId, money)
					if err != nil {
						log.Err(err).Stack().Msg("ошибка обновления баланса")
					}
				}
			}
		}
	}
}

func (r Repository) InitClientsQueues() {
	for {
		strId, err := r.rdb.LPop(context.Background(), "active_clients").Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				break
			}
			log.Err(err).Stack().Msg("ошибка получения данных из Redis")
		}
		userId, err := uuid.Parse(strId)
		if err != nil {
			log.Err(err).Stack().Msg("ошибка Redis")
		}
		newChannel := make(chan int)
		clientsChannels[userId] = newChannel
		go r.clientsQueue(userId)
	}
}

func (r Repository) PushTask(userId uuid.UUID, operation string, money int) error {
	_, ok := clientsChannels[userId]
	if !ok {
		newChannel := make(chan int)
		clientsChannels[userId] = newChannel

		_, err := r.rdb.RPush(context.Background(), "active_clients", userId.String()).Result()
		if err != nil {
			log.Err(err).Stack().Msg("ошибка записи данных в Redis")
		}

		go r.clientsQueue(userId)
	}
	_, err := r.rdb.RPush(context.Background(), userId.String(), fmt.Sprintf("%s%d", operation, money)).Result()
	if err != nil {
		log.Err(err).Stack().Msg("ошибка записи данных в Redis")
	}
	return nil
}
