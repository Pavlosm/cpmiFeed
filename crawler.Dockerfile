FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY cmd/crawlerApp ./
COPY pkg ./pkg
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /crawlerApp

ENV MONGO_ADDRESS=localhost:27017
ENV MONGO_DATABASE=cpmiFeed
ENV MONGO_PASSWORD=example
ENV MONGO_USERNAME=root
ENV KAFKA_BROKERS="broker1:9091,broker2:9092,broker3:9093"
ENV KAFKA_CONSUMER_GROUP_ID=cpmiEventsConsumer
ENV KAFKA_EVENTS_TOPIC=cpmiEvents

CMD ["/crawlerApp"]