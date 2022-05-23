package main

import (
	"TransactionsQueues/internal"
	"TransactionsQueues/internal/handler"
	"TransactionsQueues/internal/repository"
	"TransactionsQueues/pkg/pg_client"
	"TransactionsQueues/pkg/redis_client"
	"context"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// ZEROLOG
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// DOTENV
	if err := godotenv.Load(); err != nil {
		log.Fatal().Stack().Err(err).Msg("ошибка загрузки .env файла")
		return
	}

	// REDIS
	rdb, err := redis_client.NewRedisDB(redis_client.ConfigRedis{
		Host:     os.Getenv("REDIS_ADDRESS"),
		Port:     os.Getenv("REDIS_EXT_PORT"),
		Password: "",
		DB:       0,
	})
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("ошибка подключения к Redis")
		return
	}

	// POSTGRES
	db, err := pg_client.NewPostgresDB(pg_client.ConfigPG{
		Host:     os.Getenv("DB_POSTGRES_HOST"),
		Port:     os.Getenv("DB_POSTGRES_EXT_PORT"),
		Username: os.Getenv("DB_POSTGRES_USER"),
		Password: os.Getenv("DB_POSTGRES_PASSWORD"),
		DBName:   os.Getenv("DB_POSTGRES_DB"),
		SSLMode:  "disable",
		TimeZone: "Europe/Moscow",
	})
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("ошибка подключения к базе данных")
	}

	// Init
	repos := repository.NewRepository(db, rdb)
	handlers := handler.NewHandler(repos)
	srv := new(internal.Server)
	go func() {
		if err := srv.Run(os.Getenv("API_GOLANG_PORT"), handlers.InitRoutes()); err != nil {
			log.Fatal().Stack().Err(err).Msg("ошибка запуска http сервера")
		}
	}()

	go repos.InitClientsQueues()

	log.Info().Msg("Сервер запущен!..")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info().Msg("Сервер остановлен!..")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal().Stack().Err(err).Msg("ошибка остановки сервера")
	}
}
