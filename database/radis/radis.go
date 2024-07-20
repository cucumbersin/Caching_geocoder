package radis

import (
	"log"
	"log/slog"
	"strconv"

	"github.com/go-redis/redis"
)

type Radis struct {
	Logger      *slog.Logger
	RedisClient *redis.Client
}

func (RDB *Radis) SaveCordinat(key, x, y, apiID string) error {
	err := RDB.RedisClient.HMSet(key, map[string]interface{}{
		"x":     x,
		"y":     y,
		"apiID": apiID,
	}).Err()
	if err != nil {
		log.Fatalf("could not set hash: %v", err)
		return err
	}

	return nil
}

func (RDB *Radis) GetCoordinates(key string) (string, string, int, error) {
	result, err := RDB.RedisClient.HGetAll(key).Result()
	if err != nil {
		log.Fatalf("could not get hash: %v", err)
		return "", "", -1, err
	}
	apiID, _ := strconv.Atoi(result["apiID"])
	return result["x"], result["y"], apiID, nil
}

func CreateRadis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.5.0.3:6382", // адрес Redis сервера
		Password: "",              // пароль, если есть
		DB:       0,               // номер базы данных
	})
	return rdb
}
