FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /shepherd cmd/main.go

EXPOSE 8080

CMD ["/shepherd"]
