# Eth Blocks Service

The project is a service for fetching and preserving Ethereum blocks and their transactions

The service includes:
* 8 microservices
* 6 RabbitMQ instances for inter-service communication
* 4 databases: 2 write-replicas and 2 read-replicas
* Redis as a caching layer

Regarding the microservices:
* The consumer-facing microservice is built on Fiber, which utilizes Fasthttp, the fastest HTTP engine for Go
* The remaining microservices are designed with minimal dependencies, without relying on any frameworks

Additional Features:
* Distributed Transactions: the microservices are capable of handling transactions that involve multiple network hosts
* Cron Job: this automated task is responsible for fetching the latest blocks and transactions.
* Repository Pattern: all microservices employ the repository pattern


Infrastructure:
* Dockerization: The entire system, comprised of 20 elements, is containerized
* Makefile: The entire system can be launched with a single command
* Live-Reload: Each microservice automatically reloads upon any changes

The entire solution is characterized by being:
* Scalable
* Highly decoupled
* Fault-tolerant
* Technology agnostic
* Reusable
* Secure
* Optimized for the Cloud

Bonus:
* The front-end is built using Rust and WebAssembly



![Architecture Diagram](./assets/diagrams/eth-block-service.drawio-full.png)




