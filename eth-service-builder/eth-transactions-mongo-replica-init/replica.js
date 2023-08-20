rs.initiate({
    _id : "ethTransactionsMongoReplica",
    members: [
        { _id: 0, host: "eth-transactions-mongo:27017" },
        { _id: 1, host: "eth-transactions-mongo-replica:27017" }
    ]
});
