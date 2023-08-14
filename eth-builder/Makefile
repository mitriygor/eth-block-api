clean:
	@find ./../ -name .tmp -type d -exec rm -rf {} +

purge:
	@find ./../ -name .db-data -type d -exec rm -rf {} +

build: clean
	@docker-compose build --no-cache

run: build
	@docker-compose up

rerun:
	@make purge
	@make build
	@docker-compose up  -d
	@echo "Rerunning, please wait..."
	@for i in {1..30}; do \
		for c in / - \\ \|; do \
			printf "\r%c" "$$c"; \
			sleep 0.25; \
		done; \
	done; printf "\n"
	@docker exec -it eth-builder-mongo-1 mongo -u admin -p password123 --authenticationDatabase admin /docker-entrypoint-initdb.d/replica.js
	@echo "Done rerunning!"

stop:
	@docker-compose down
