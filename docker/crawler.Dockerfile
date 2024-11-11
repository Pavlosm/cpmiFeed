FROM golang:1.23-bullseye AS build

# Create a non-root user for running the application
RUN useradd -u 1001 nonroot

WORKDIR /app

COPY go.mod go.sum ./
COPY cmd/crawlerApp ./
COPY pkg ./pkg
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/crawlerApp

FROM scratch

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/crawlerApp /crawlerApp

USER nonroot

CMD ["/crawlerApp"]