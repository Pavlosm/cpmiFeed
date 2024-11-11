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

ARG DEFAULT_PORT=8099
ENV WEB_APP_PORT=$DEFAULT_PORT

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/webApp /webApp

USER nonroot

EXPOSE $WEB_APP_PORT

CMD ["/webApp"]