resource "scaleway_account_project" "main" {
  name = var.project_name

  lifecycle {
    prevent_destroy = true
  }
}

resource "scaleway_rdb_instance" "main" {
  project_id = scaleway_account_project.main.id

  name = "${var.project_name}-db"

  node_type     = var.db_node_type
  is_ha_cluster = false

  volume_type       = "bssd"
  volume_size_in_gb = 10

  engine = "PostgreSQL-${var.db_postgres_version}"

  user_name = var.db_admin_user_name
  password  = var.db_admin_password

  tags = [
    var.project_name,
    "db",
    "production"
  ]
}

// Get the IP of the current machine
data "http" "icanhazip" {
  url = "https://ipv4.icanhazip.com"
}

// With Scaleway Serverless Containers, we do not have a static IP
// Instead, we allow traffic from all IPs
// See: https://as12876.net/
resource "scaleway_rdb_acl" "main" {
  instance_id = scaleway_rdb_instance.main.id

  dynamic "acl_rules" {
    for_each = [
      "62.210.0.0/16",
      "195.154.0.0/16",
      "212.129.0.0/18",
      "62.4.0.0/19",
      "212.83.128.0/19",
      "212.83.160.0/19",
      "212.47.224.0/19",
      "163.172.0.0/16",
      "51.15.0.0/16",
      "151.115.0.0/16",
      "51.158.0.0/15",
    ]
    content {
      ip          = acl_rules.value
      description = "Allow Scaleway IPs"
    }
  }

  acl_rules {
    ip          = "${trimspace(data.http.icanhazip.response_body)}/32"
    description = "Allow current IP"
  }
}

resource "scaleway_rdb_database" "main" {
  instance_id = scaleway_rdb_instance.main.id
  name        = var.db_name
}

resource "scaleway_rdb_privilege" "admin" {
  instance_id   = scaleway_rdb_instance.main.id
  user_name     = var.db_admin_user_name
  database_name = scaleway_rdb_database.main.name

  permission = "all"
}

resource "scaleway_rdb_user" "main" {
  instance_id = scaleway_rdb_instance.main.id

  name     = var.db_user_name
  password = var.db_password
}

resource "scaleway_rdb_privilege" "main" {
  instance_id   = scaleway_rdb_instance.main.id
  user_name     = scaleway_rdb_user.main.name
  database_name = scaleway_rdb_database.main.name

  // By default, is unable to run migrations
  // With Atlas, we run the migrations separately 
  // so it's not a problem
  permission = "readwrite"
}

locals {
  db_endpoint = scaleway_rdb_instance.main.load_balancer[0]
}

resource "scaleway_container_namespace" "main" {
  project_id = scaleway_account_project.main.id

  name = var.project_name

  environment_variables = {
    "ENVIRONMENT" = "production",
  }

  secret_environment_variables = {
    "PG_CONN" = "postgres://${scaleway_rdb_user.main.name}:${scaleway_rdb_user.main.password}@${local.db_endpoint.ip}:${local.db_endpoint.port}/${var.db_name}"
  }
}

resource "scaleway_container" "api" {
  namespace_id = scaleway_container_namespace.main.id

  name           = "${var.project_name}-api"
  description    = "API for ${var.project_name}"
  registry_image = local.api_image
  port           = 8080
  privacy        = "public"
  deploy         = true

  cpu_limit    = 500
  memory_limit = 256

  environment_variables = {
    "FRONTEND_URL" = "https://${var.project_name}.com",
  }

  depends_on = [
    null_resource.build_api,
  ]
}

locals {
  crawler_memory_limit = 1024
}

resource "scaleway_container" "crawler" {
  namespace_id = scaleway_container_namespace.main.id

  name           = "${var.project_name}-crawler"
  description    = "Crawler for ${var.project_name}"
  registry_image = local.crawler_image
  port           = 8081
  privacy        = "private"
  deploy         = true

  cpu_limit    = 2000
  memory_limit = local.crawler_memory_limit

  max_scale = 1

  environment_variables = {
    "CRAWL_DURATION" = "15m",
    "CRAWLER_COUNT"  = "3",
    "GOMEMLIMIT"     = "${floor(0.85 * local.crawler_memory_limit)}MiB",
    "GOGC"           = "100"
  }

  timeout = 900 // 15 minutes

  depends_on = [
    null_resource.build_crawler,
  ]
}

resource "scaleway_container_cron" "crawl_every_hour" {
  container_id = scaleway_container.crawler.id

  // Every hour, at 17 minutes
  // (use a weird time to avoid collisions with other cron jobs)
  schedule = "17 * * * *"

  // Body is ignored, but we need to provide something
  args = "{}"
}
