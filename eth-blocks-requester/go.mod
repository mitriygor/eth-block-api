module eth-blocks-requester

go 1.19

require (
	eth-helpers v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/rabbitmq/amqp091-go v1.8.1
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
)

replace eth-helpers v0.0.0 => ../eth-helpers
