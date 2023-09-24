provider "docker" {
  host = "unix:///var/run/docker.sock"

  registry_auth {
    address  = scaleway_container_namespace.main.registry_endpoint
    username = "nologin"
    password = scaleway_iam_api_key.registry_push.secret_key
  }
}

locals {
  api_image     = "${scaleway_container_namespace.main.registry_endpoint}/${var.project_name}-api:${var.docker_image_tag}"
  crawler_image = "${scaleway_container_namespace.main.registry_endpoint}/${var.project_name}-crawler:${var.docker_image_tag}"
}

// While it's possible to build the image via the Docker provider,
// it doesn't support BUILDKIT, see: https://github.com/docker/for-linux/issues/1136
// So we use a local-exec provisioner to build the image and push it to the registry
resource "null_resource" "build_api" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command = "docker build --target api --tag ${local.api_image} .."
    environment = {
      // In case the user has set DOCKER_BUILDKIT=0, we override it
      "DOCKER_BUILDKIT" = "1"
    }
  }

  provisioner "local-exec" {
    command = "docker push ${local.api_image}"
  }
}

resource "null_resource" "build_crawler" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command = "docker build --target crawler --tag ${local.crawler_image} .."
    environment = {
      // In case the user has set DOCKER_BUILDKIT=0, we override it
      "DOCKER_BUILDKIT" = "1"
    }
  }

  provisioner "local-exec" {
    command = "docker push ${local.crawler_image}"
  }
}
