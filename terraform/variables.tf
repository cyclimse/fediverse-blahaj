variable "project_name" {
  type    = string
  default = "blahaj"
}

variable "db_node_type" {
  type    = string
  default = "db-play2-nano"
}

variable "db_postgres_version" {
  type    = number
  default = 15
}

variable "db_admin_user_name" {
  type    = string
  default = "admin"
}

variable "db_admin_password" {
  type      = string
  sensitive = true
}

variable "db_name" {
  type    = string
  default = "blahaj"
}

variable "db_user_name" {
  type    = string
  default = "blahaj"
}

variable "db_password" {
  type      = string
  sensitive = true
}

variable "docker_image_tag" {
  type    = string
  default = "latest"
}
