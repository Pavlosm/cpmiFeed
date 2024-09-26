build:
	@echo "Building all applications..."
	@echo "Building crawlerApp..."
	@cd cmd/crawlerApp && go build -o ../../bin/crawlerApp
	@echo "crawlerApp built successfully."
	@echo "Building eventConsumer..."
	@cd cmd/eventConsumer && go build -o ../../bin/eventConsumer
	@echo "eventConsumer built successfully."
	@echo "Building webApp..."
	@cd cmd/webApp && go build -o ../../bin/webApp
	@echo "webApp built successfully."
	@echo "All applications built successfully."

start-crawler:
	@echo "Starting crawlerApp..."
	@cd cmd/crawlerApp && go run .
	@echo "crawlerApp running successfully."

stop-crawler:
	@echo "Stopping crawlerApp..."
	@cd bin && pkill -f "crawlerApp"
	@echo "crawlerApp stopped successfully."

start-eventConsumer:
	@echo "Starting eventConsumer..."
	@cd cmd/eventConsumer && go run .
	@echo "eventConsumer running successfully."

stop-eventConsumer:
	@echo "Stopping eventConsumer..."
	@cd bin && pkill -f "eventConsumer"
	@echo "eventConsumer stopped successfully."

start-webApp:
	@echo "Starting webApp..."
	@cd cmd/webApp && go run .
	@echo "webApp running successfully."

stop-webApp:
	@echo "Stopping webApp..."
	@cd bin && pkill -f "webApp"
	@echo "webApp stopped successfully."

stop-all:
	@echo "Stopping all applications..."
	@echo "Stopping crawlerApp..."
	@cd bin && pkill -f "crawlerApp"
	@echo "crawlerApp stopped successfully."
	@echo "Stopping eventConsumer..."
	@cd bin && pkill -f "eventConsumer"
	@echo "eventConsumer stopped successfully."
	@echo "Stopping webApp..."
	@cd bin && pkill -f "webApp"
	@echo "webApp stopped successfully."

clean: 
	@cd bin && rm -f crawlerApp eventConsumer webApp