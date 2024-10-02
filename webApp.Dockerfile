FROM golang:1.23-bullseye AS build

# Create a non-root user for running the application
RUN useradd -u 1001 nonroot

WORKDIR /app

COPY go.mod go.sum ./
COPY cmd/webApp/ ./cmd/webApp
COPY pkg ./pkg
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/webApp ./cmd/webApp

FROM scratch

ENV MONGO_ADDRESS=mongo:27017
ENV MONGO_DATABASE=cpmiFeed
ENV MONGO_PASSWORD=example
ENV MONGO_USERNAME=root
ENV KAFKA_BROKERS="broker1:9091,broker2:9092,broker3:9093"
ENV KAFKA_CONSUMER_GROUP_ID=cpmiEventsConsumer
ENV KAFKA_EVENTS_TOPIC=cpmiEvents
ENV WEB_APP_PORT=8099

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/webApp /webApp

USER nonroot

EXPOSE 8099

CMD ["/webApp"]