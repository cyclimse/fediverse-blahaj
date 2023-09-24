output "db_admin_connection_string" {
  value       = "postgresql://${var.db_admin_user_name}:${var.db_admin_password}@${local.db_endpoint.ip}:${local.db_endpoint.port}/${var.db_name}"
  sensitive   = true
  description = "Connection string for the database admin user."
}
