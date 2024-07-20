FROM golang

WORKDIR /app

COPY . .

RUN go build cmd/Caching_geocoder/main.go

CMD ["./main"]