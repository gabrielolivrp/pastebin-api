variable "project_name" {
  type        = string
  description = "Project name"
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources"
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
