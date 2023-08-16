rs.initiate({
    _id : "mongoReplicaSet2",
    members: [
        { _id: 0, host: "mongo:27017" },
        { _id: 1, host: "mongo-for-eth-blocks-reading:27017" }
    ]
});
