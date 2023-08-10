# The API Builder

The builder provides capabilities to run the multi-service application. It consists of API, Emitter, and Listener services and uses Redis, MongoDB, and RabbitMQ.

### Prerequisites

- Docker and Docker Compose installed
- Make sure the required ports are available on the system

### Installation

1. Clone the repository to the local machine
2. Navigate to the project directory
3. Run the application using the Makefile commands

## Usage

### Makefile Commands

- **clean**: Remove the .tmp directory
   ```bash
   make clean

- **build**: Build the Docker images without using cache
   ```bash
   make build

- **run**: Build and run the Docker containers
   ```bash
   make run

- **stop**: Stop the running Docker containers
   ```bash
   make stop
  

The project uses the [Air](https://github.com/cosmtrek/air) package for live reloading. It watches for file changes and automatically restarts the application. The Air configuration is stored in the .air.toml file.