FROM golang:1.23.6

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o banking-service ./cmd/main.go

EXPOSE 8080

CMD ["./banking-service"]
