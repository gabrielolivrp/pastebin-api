variable "project_name" {
  type        = string
  description = "Project name"
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources"
}

variable "db_username" {
  type        = string
  description = "The username for the database"
}

variable "db_password" {
  type        = string
  description = "The password for the database"
}

variable "db_name" {
  type        = string
  description = "The name of the database"
}

variable "cluster_security_group_id" {
  type        = string
  description = "The security group ID of the EKS cluster"
}

variable "vpc_id" {
  type        = string
  description = "The VPC ID where the RDS instance will be deployed"
}

variable "private_subnet_1a" {
  type        = string
  description = "The ID of the private subnet in availability zone 1a"
}

variable "private_subnet_1b" {
  type        = string
  description = "The ID of the private subnet in availability zone 1b"
}
