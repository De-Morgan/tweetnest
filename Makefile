


## up: starts all containers in the background
up: 
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"


## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down  --rmi 'local'
	@echo "Done!"