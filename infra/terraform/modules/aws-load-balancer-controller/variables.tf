variable "project_name" {
  type        = string
  description = "Project name"
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources"
}

variable "oidc" {
  type        = string
  description = "The OIDC provider URL for the EKS cluster"
}

variable "cluster_name" {
  type        = string
  description = "The name of the EKS cluster"
}

variable "vpc_id" {
  type        = string
  description = "The VPC ID where the EKS cluster is deployed"
}
