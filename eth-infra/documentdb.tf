provider "aws" {
  region = "us-east-1"  # Change to your preferred region
}

resource "aws_docdb_cluster" "my_docdb_cluster" {
  cluster_identifier      = "my-docdb-cluster"
  engine                   = "docdb"
  engine_version           = "5.0.0"  # Adjust to your preferred version (Note: At the time of my last update in September 2021, the latest version was 4.0.0, so check the latest version)
  master_username          = "mydbadmin"
  master_password          = "mysecurepassword123"
  backup_retention_period  = 5
  preferred_backup_window  = "07:00-09:00"
  skip_final_snapshot      = true
}

resource "aws_docdb_cluster_instance" "my_docdb_instance" {
  identifier              = "my-docdb-instance"
  cluster_identifier      = aws_docdb_cluster.my_docdb_cluster.id
  instance_class          = "db.t3.medium"
  engine                  = "docdb"
  auto_minor_version_upgrade = true
}

output "endpoint" {
  value = aws_docdb_cluster.my_docdb_cluster.endpoint
}

output "reader_endpoint" {
  value = aws_docdb_cluster.my_docdb_cluster.reader_endpoint
}
