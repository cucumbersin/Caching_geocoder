package database

import (
	"log/slog"

	"example.com/m/database/radis"
	"github.com/go-redis/redis"
)

type CoordinatesDb interface {
	SaveCordinat(key, x, y, apiID string) error
	GetCoordinates(key string) (string, string, int, error)
}

func CreateRadis(Logger *slog.Logger, addr string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.5.0.3:6379", // адрес Redis сервера
		Password: "",              // пароль, если есть
		DB:       0,               // номер базы данных
	})
	el := &radis.Radis{
		Logger:      Logger,
		RedisClient: rdb,
	}
	DB = el
}

var DB CoordinatesDb
