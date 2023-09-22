# Ethereum Blocks Service

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
* [Dockerization](https://github.com/mitriygor/eth-block-api/tree/main/eth-service-builder): The entire system, comprised of 20 elements, is containerized
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



![Architecture Diagram](./assets/diagrams/eth-block-service-8.png)


## Endpoints

Below, we detail the available GET endpoints and the type of information they return.

### {{URL}}/eth-blocks/0x110cfed
Returns information about the Ethereum block identified by the hex number **0x110cfed**.

- **Endpoint:** {{URL}}/eth-blocks/0x110cfed
- **Method:** GET
- **URL Params:** None
- **Data Params:** None
- **Success Response:**
  - **Code:** 200
  - **Content:** JSON object containing detailed block information
- **Error Response:**
  - **Code:** 404
  - **Content:** {"error": "Block not found"}


### {{URL}}/eth-blocks/0x22c5789c9c2daf92c206489ad61cd20304084e17e4a6b215e194750b67e17996
Returns information about the Ethereum block identified by the hash **0x22c5789c9c2daf92c206489ad61cd20304084e17e4a6b215e194750b67e17996**.

- **Endpoint:** {{URL}}/eth-blocks/0x22c5789c9c2daf92c206489ad61cd20304084e17e4a6b215e194750b67e17996
- **Method:** GET
- **URL Params:** None
- **Data Params:** None
- **Success Response:**
  - **Code:** 200
  - **Content:** JSON object containing detailed block information
- **Error Response:**
  - **Code:** 404
  - **Content:** {"error": "Block not found"}

### {{URL}}/eth-transactions/0x5a53ff76232b1fdc722583e0afd4f62a70dec6ae8e52347958a94b2957156144
Returns information about the Ethereum transaction identified by the hash **0x5a53ff76232b1fdc722583e0afd4f62a70dec6ae8e52347958a94b2957156144**.

- **Endpoint:** {{URL}}/eth-transactions/0x5a53ff76232b1fdc722583e0afd4f62a70dec6ae8e52347958a94b2957156144
- **Method:** GET
- **URL Params:** None
- **Data Params:** None
- **Success Response:**
  - **Code:** 200
  - **Content:** JSON object containing detailed transaction information
- **Error Response:**
  - **Code:** 404
  - **Content:** {"error": "Transaction not found"}

### {{URL}}/eth-events/0x2cc846fff0b08fb3bffad71f53a60b4b6e6d6482
Returns information about Ethereum events associated with the gas used identified by the hex **0x2cc846fff0b08fb3bffad71f53a60b4b6e6d6482**.

- **Endpoint:** {{URL}}/eth-events/0x2cc846fff0b08fb3bffad71f53a60b4b6e6d6482
- **Method:** GET
- **URL Params:** None
- **Data Params:** None
- **Success Response:**
  - **Code:** 200
  - **Content:** JSON object containing detailed events information
- **Error Response:**
  - **Code:** 404
  - **Content:** {"error": "Events not found"}



