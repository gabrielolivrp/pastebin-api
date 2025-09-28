variable "project_name" {
  type        = string
  description = "Project name"
}

variable "tags" {
  type        = map(string)
  description = "A map of tags to add to all resources"
}

variable "public_subnet_1a" {
  type        = string
  description = "The ID of the public subnet in availability zone 1a"
}

variable "public_subnet_1b" {
  type        = string
  description = "The ID of the public subnet in availability zone 1b"
}
