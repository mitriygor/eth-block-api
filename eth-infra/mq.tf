provider "aws" {
  region = "us-east-1"  # Change to your preferred region
}

resource "aws_mq_broker" "my_broker" {
  broker_name   = "MyBroker"
  engine_type   = "RabbitMQ"

  engine_version = "3.8.22"  # Adjust to your preferred version

  host_instance_type = "mq.t3.micro"

  deployment_mode = "SINGLE_INSTANCE"

  user {
    username = "Admin"    # Change to your preferred username
    password = "MySecurePassword123!" # Change to your preferred password
  }

  # ... (other configurations as necessary)
}
