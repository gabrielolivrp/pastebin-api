output "rds_endpoint" {
  description = "RDS instance endpoint"
  value       = module.rds.rds_endpoint
}

output "rds_port" {
  description = "RDS instance port"
  value       = module.rds.rds_port
}

output "rds_database_name" {
  description = "RDS database name"
  value       = module.rds.rds_database_name
}

output "rds_username" {
  description = "RDS master username"
  value       = module.rds.rds_username
  sensitive   = true
}

output "redis_endpoint" {
  description = "Redis endpoint"
  value       = module.elasticache.redis_endpoint
}
