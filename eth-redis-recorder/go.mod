module eth-redis-recorder

go 1.19

require (
	eth-helpers v0.0.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/joho/godotenv v1.5.1
	github.com/rabbitmq/amqp091-go v1.8.1
	go.uber.org/zap v1.26.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
)

replace eth-helpers v0.0.0 => ../eth-helpers
