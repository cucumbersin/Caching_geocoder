package rest

import (
	"log/slog"
	"net/http"

	"example.com/m/internal/servic"
)

func Getcoordinates(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		address := r.URL.Query().Get("address")
		log.Debug("Getcoordinates call", "address", address)

		if address == "" {
			http.Error(w, "Missing 'address' parameter", http.StatusBadRequest)
			return
		}

		x, y := servic.Getcoordinates(address)
		w.Write([]byte("latitude=" + x + " longitude=" + y))
		log.Debug("Getcoordinates end", "latitude", x, "longitude", y)
	}
}
