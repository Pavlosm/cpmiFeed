# Build

build: build-crawlerApp build-eventConsumer build-webApp

build-crawlerApp:
	@echo Building crawlerApp...
	@cd cmd/crawlerApp && go fmt && go vet && go build -o ../../bin/crawlerApp
	@echo crawlerApp built successfully.

build-eventConsumer:
	@echo Building eventConsumer...
	@cd cmd/eventConsumer && go fmt && go vet && go build -o ../../bin/eventConsumer
	@echo eventConsumer built successfully.

build-webApp:
	@echo Building webApp...
	@cd cmd/webApp && go fmt && go vet && go build -o ../../bin/webApp
	@echo webApp built successfully.

clean: 
	@cd bin && rm -f crawlerApp eventConsumer webApp kafkaChangeStream

# Run

start-all:
	start-crawlerApp
	start-eventConsumer
	start-webApp

stop-all:
	stop-crawlerApp
	stop-eventConsumer
	stop-webApp

start-crawlerApp:
	@echo Starting crawlerApp...
	@cd cmd/crawlerApp && go run .
	@echo crawlerApp running successfully.

stop-crawlerApp:
	@echo Stopping crawlerApp...
	@cd bin && pkill -f "crawlerApp"
	@echo crawlerApp stopped successfully.

start-eventConsumer:
	@echo Starting eventConsumer...
	@cd cmd/eventConsumer && go run .
	@echo eventConsumer running successfully.

stop-eventConsumer:
	@echo Stopping eventConsumer...
	@cd bin && pkill -f "eventConsumer"
	@echo eventConsumer stopped successfully.

start-webApp:
	@echo Starting webApp...
	@cd cmd/webApp && go run .
	@echo webApp running successfully.

stop-webApp:
	@echo Stopping webApp...
	@cd bin && pkill -f "webApp"
	@echo webApp stopped successfully.

# Docker build

docker-build-crawlerApp:
	@echo Building crawlerApp Docker image...
	@docker build -t crawler -f docker/crawler.Dockerfile .
	@echo crawlerApp Docker image built successfully.

docker-build-eventConsumer:
	@echo Building eventConsumer Docker image...
	@docker build -t eventconsumer -f "docker/eventConsumer.Dockerfile" .
	@echo eventConsumer Docker image built successfully.

docker-build-webApp:
	@echo Building webApp Docker image...
	@docker build -t webapp -f docker/webApp.Dockerfile .
	@echo webApp Docker image built successfully.

# Docker run

docker-run-crawler:
	@echo Running crawlerApp Docker container...
	@docker run -d --name crawler --network=kafka_default crawler 
	@echo crawlerApp Docker container running successfully.

docker-run-eventConsumer:
	@echo "Running eventConsumer Docker container..."
	@docker run -d --name eventconsumer --network=kafka_default --network=localsetup_default eventconsumer
	@echo "eventConsumer Docker container running successfully."

docker-run-webApp:
	@echo Running webApp Docker container...
	@docker run -d --name webapp --network=localsetup_default -p 8095:8095 webapp
	@echo webApp Docker container running successfully.