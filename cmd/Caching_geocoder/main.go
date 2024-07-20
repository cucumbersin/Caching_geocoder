package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"example.com/m/database"
	"example.com/m/internal/config"
	"example.com/m/internal/servic"
	middleware "example.com/m/internal/transport/middlewate"
	"example.com/m/internal/transport/rest"
	"github.com/ekomobile/dadata/v2/client"
	"github.com/gorilla/mux"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	servic.YandexApiKey = cfg.API.YandexKey
	Creds := client.Credentials{
		ApiKeyValue:    cfg.API.DaDataApi.ApiKeyValue,
		SecretKeyValue: cfg.API.DaDataApi.SecretKeyValue,
	}
	servic.Creds = Creds
	servic.GeocodeMapsCo_key = cfg.API.GeocodeMapsCoKey

	log.Debug("init log")
	log.Debug("", "cfg.Radis", cfg.Radis)
	router := mux.NewRouter()
	logMiddleware := middleware.LoggMiddleware{Logger: log}
	router.Use(logMiddleware.Middleware)
	database.CreateRadis(log, cfg.DB.Radis)
	router.HandleFunc("/", rest.Getcoordinates(log)).Methods("GET")

	servic.Log = log

	serv := &http.Server{
		Handler:      router,
		Addr:         cfg.HTTP_server.IP + ":" + strconv.Itoa(cfg.HTTP_server.Port),
		WriteTimeout: cfg.HTTP_server.Timeout,
		ReadTimeout:  cfg.HTTP_server.Timeout,
		IdleTimeout:  cfg.HTTP_server.Iddle_timeout,
	}
	log.Error("server stop %s", serv.ListenAndServe())
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
