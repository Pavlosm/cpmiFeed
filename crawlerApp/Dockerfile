FROM golang:1.20-alpine AS base

WORKDIR /app

COPY ../go.mod ./
COPY ../go.sum ./
RUN go mod download

COPY *.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o /crawlerApp

CMD ["/crawlerApp"]