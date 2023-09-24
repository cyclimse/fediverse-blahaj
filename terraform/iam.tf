resource "scaleway_iam_application" "registry_push" {
  name = "${var.project_name}-registry-push"
}

resource "scaleway_iam_policy" "registry_full_access" {
  name           = "${var.project_name}-registry-full-access"
  description    = "Give full access to container registry."
  application_id = scaleway_iam_application.registry_push.id
  rule {
    project_ids          = [scaleway_account_project.main.id]
    permission_set_names = ["ContainerRegistryFullAccess"]
  }
}

resource "scaleway_iam_api_key" "registry_push" {
  application_id = scaleway_iam_application.registry_push.id
  description    = "Ephemeral API key to push to registry."

  expires_at = timeadd(timestamp(), "10m")
}
