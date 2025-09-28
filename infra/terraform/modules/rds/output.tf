output "rds_endpoint" {
  description = "RDS instance endpoint"
  value       = aws_db_instance.postgres_instance.endpoint
  sensitive   = false
}

output "rds_port" {
  description = "RDS instance port"
  value       = aws_db_instance.postgres_instance.port
}

output "rds_database_name" {
  description = "RDS database name"
  value       = aws_db_instance.postgres_instance.db_name
}

output "rds_username" {
  description = "RDS master username"
  value       = aws_db_instance.postgres_instance.username
  sensitive   = true
}
