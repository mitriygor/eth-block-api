# Makefile for handling common development tasks using make command

# The 'clean' target deletes all .tmp directories in the parent directory.
clean:
	@# Using 'find' command to search for directories named '.tmp' and delete them.
	@# './../' - denotes searching should begin from the parent directory.
	@# '-name .tmp' - matches directories with the name '.tmp'.
	@# '-type d' - ensures only directories (not files) are found.
	@# '-exec rm -rf {} +' - deletes (without asking for confirmation) each found directory.
	@find ./../ -name .tmp -type d -exec rm -rf {} +

# The 'purge' target deletes all .db-data directories in the parent directory.
purge:
	@# Similar to 'clean', but matches directories with the name '.db-data'.
	@find ./../ -name .db-data -type d -exec rm -rf {} +

# The 'build' target first ensures a clean environment and then builds Docker containers.
build: clean
	@# 'docker-compose build --no-cache' builds Docker images without using any cached layers.
	@# '--no-cache' forces Docker to build the images from scratch, ensuring the newest code changes are applied.
	@docker-compose build --no-cache

# The 'run' target builds (ensuring a clean environment before) and then starts Docker containers.
run: build
	@# 'docker-compose up' starts up the Docker containers based on the Docker Compose configuration.
	@docker-compose up

# The 'run-fe' target navigates to the 'eth-fe' directory and runs the 'trunk serve' command.
run-fe:
	@# Using 'cd' to navigate to the 'eth-fe' directory relative to the Makefile.
	@# '$(MAKE) -C [DIR]' is used to run make targets in a specified directory ([DIR]).
	@# 'trunk serve' is assumed to start a development server for a Rust Trunk application.
	cd ../eth-fe && trunk serve

# The 'run-all' target runs both the main application and the frontend.
run-all:
	@# 'make run' - triggers the 'run' target which is assumed to build and start the main application.
	make run
	@# 'make run-fe' - triggers the 'eth-fe' target which is assumed to run the frontend part of the application.
	make run-fe

# The 'rerun' target purges, builds (with a clean before), and then runs Docker containers.
rerun:
	@# First it purges any data by calling the 'purge' target.
	@make purge

	@# Then it builds the Docker images ensuring a clean environment.
	@make build

	@# Lastly, it starts up the Docker containers.
	@docker-compose up

# The 'stop' target stops and removes Docker containers, networks, volumes, and images defined in the Docker Compose configuration.
stop:
	@# 'docker-compose down' stops and removes containers, as well as associated volumes and networks.
	@docker-compose down

# The 'test' target runs all Golang tests in the specified projects.
test:
	@# The 'go test' command runs tests in the Golang projects.
	@# '-count=1' prevents test caching to ensure all tests are always run.
	@# '-v' enables verbose mode, logging all tests as they are run.
	@# Each line below navigates into a service directory, runs the Go tests, and then returns to the original directory.
	@cd ../eth-api && go test -count=1 -v ./... && cd -
	@cd ../eth-blocks-recorder && go test -count=1 -v ./... && cd -
	@cd ../eth-blocks-requester && go test -count=1 -v ./... && cd -
	@cd ../eth-blocks-scheduler && go test -count=1 -v ./... && cd -
	@cd ../eth-helpers && go test -count=1 -v ./... && cd -
	@cd ../eth-redis-recorder && go test -count=1 -v ./... && cd -
	@cd ../eth-transactions-recorder && go test -count=1 -v ./... && cd -
	@cd ../eth-transactions-requester && go test -count=1 -v ./... && cd -
	@cd ../eth-transactions-scheduler && go test -count=1 -v ./... && cd -
