module eth-transactions-requester

go 1.19

require (
	github.com/joho/godotenv v1.5.1
	github.com/rabbitmq/amqp091-go v1.8.1
	eth-helpers v0.0.0
)

replace eth-helpers v0.0.0 => ../eth-helpers