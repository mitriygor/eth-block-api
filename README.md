# Eth Blocks API

![Architecture Diagram](./assets/diagrams/eth-block-service.drawio.png)


### Blocks Fetching

1. Checking the current latest block saved in the storage
2. Getting number of the latest block from the ETH RPC, and based the difference with the current one, fetching lacking latest blocks
3. Pushing the fetched blocks to the Queue
4. Listener received the blocks from the Queue
5. Listener saves the blocks to the storage 
6. Setting the block number to the memory


### Blocks Requesting

1. Sending a request to the API
2. API fetching the block from the storage