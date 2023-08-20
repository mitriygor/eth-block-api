rs.initiate({
    _id : "ethBlocksMongoReplica",
    members: [
        { _id: 0, host: "eth-blocks-mongo:27017" },
        { _id: 1, host: "eth-blocks-mongo-replica:27017" }
    ]
});
