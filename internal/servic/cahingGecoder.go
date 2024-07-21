package servic

import (
	"context"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"example.com/m/database"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
)

var Log *slog.Logger

func Getcoordinates(addres string) (string, string) {
	Log.Debug("servic Getcoordinates call", "addres", addres)
	x, y, apiID, err := database.DB.GetCoordinates(addres)
	if err != nil {
		Log.Error("GetCoordinates err", "err", err)
	}
	Log.Debug("servic Getcoordinates call", "x", x, "y", y, "apiID", apiID)
	if apiID == 0 {
		Log.Debug("value dont exist in db")
		x, y = yandexApi(addres)
		if !(x == "" || y == "") {
			database.DB.SaveCordinat(addres, x, y, "1")
			return x, y
		}
		x, y = daDataApi(addres)
		if !(x == "" || y == "") {
			database.DB.SaveCordinat(addres, x, y, "1")
			return x, y
		}
		if err != nil {
			Log.Error("SaveCordinat err", "err", err)
		}
		x, y = geocodeMapsCo(addres)
		if err != nil {
			Log.Error("SaveCordinat err", "err", err)
		}
	}
	return x, y
}

// var YandexApiKey string = "172f06a9-b24f-4cc8-ad28-12afbdbc2a35"
var YandexApiKey string

func yandexApi(addres string) (string, string) {
	Log.Debug("servic yandex_api call", "addres", addres)
	baseURL := "https://geocode-maps.yandex.ru/1.x/?apikey="
	baseURL += YandexApiKey + "&geocode="
	encodedQuery := url.QueryEscape(addres)
	fullURL := fmt.Sprintf("%s?q=%s", baseURL, encodedQuery)
	resp, err := http.Get(fullURL)
	if err != nil {
		Log.Error("yandex_api", "err", err)
		return "", ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var str string = string(body)
	f := strings.Index(str, "<pos>")
	if f == -1 {
		return "", ""
	}
	f = f + len("<pos>")
	s := f
	for ; s < len(str); s++ {
		if str[s] == ' ' {
			break
		}
	}
	x := str[f:s]
	f = s + 1
	for ; s < len(str); s++ {
		if str[s] == '<' {
			break
		}
	}
	y := str[f:s]
	return y, x
}

var Creds client.Credentials

func daDataApi(addres string) (string, string) {
	api := dadata.NewSuggestApi(client.WithCredentialProvider(&Creds))
	Log.Debug("daDataApi", "addres", addres)
	params := suggest.RequestParams{
		Query: addres,
	}

	suggestions, err := api.Address(context.Background(), &params)
	if err != nil {
		Log.Error("daDataApi", "err", err)
		return "", ""
	}
	var x, y string
	Log.Debug("daDataApi", "len s.Value", len(suggestions))
	for _, s := range suggestions {
		Log.Debug("daDataApi", "s.Value", s.Value)
		x = s.Data.GeoLat
		y = s.Data.GeoLon
	}
	return x, y
}

// var GeocodeMapsCo_key string = "669a1ca23ac61537414478ymc268c3c"
var GeocodeMapsCo_key string

// https://geocode.maps.co/
func geocodeMapsCo(addres string) (string, string) {
	//https://geocode.maps.co/search?q=Литейный проспект, 21, Санкт-Петербург&api_key=669a1ca23ac61537414478ymc268c3c&format=xml
	Log.Debug("servic geocodeMapsCo call", "addres", addres)
	baseURL := "https://geocode.maps.co/search?api_key="
	baseURL += GeocodeMapsCo_key + "&format=xml"
	Log.Debug("servic geocodeMapsCo call", "baseURL", baseURL)
	encodedQuery := url.QueryEscape(addres)
	fullURL := fmt.Sprintf("%s&q=%s", baseURL, encodedQuery)
	Log.Debug("servic geocodeMapsCo call", "fullURL", fullURL)
	resp, err := http.Get(fullURL)
	if err != nil {
		Log.Error("yandex_api", "err", err)
		return "", ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var str string = string(body)
	fmt.Println(str)
	f := strings.Index(str, "lat=")
	if f == -1 {
		return "", ""
	}
	f = f + len("lat=") + 1
	fmt.Println("f=", f, "str[f]", str[f])
	s := f
	for ; s < len(str); s++ {
		if str[s] == '\'' {
			break
		}
	}
	x := str[f:s]
	f = s + len("' lon='")
	s = f
	for ; s < len(str); s++ {
		if str[s] == '\'' {
			break
		}
	}
	fmt.Println("f=", f, "s=", s)
	y := str[f:s]
	return x, y
}
