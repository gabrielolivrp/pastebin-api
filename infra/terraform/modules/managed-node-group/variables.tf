variable "project_name" {
  type        = string
  description = "Project name"
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources"
}

variable "cluster_name" {
  type        = string
  description = "The name of the EKS cluster to join"
}

variable "private_subnet_1a" {
  type        = string
  description = "The ID of the private subnet in availability zone 1a"
}

variable "private_subnet_1b" {
  type        = string
  description = "The ID of the private subnet in availability zone 1b"
}
