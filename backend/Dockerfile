FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /helm-deployer cmd/main.go

EXPOSE 8080

CMD ["/helm-deployer"]
