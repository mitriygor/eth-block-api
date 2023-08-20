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
	@docker-compose up

stop:
	@docker-compose down
