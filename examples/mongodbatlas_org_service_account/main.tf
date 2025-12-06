resource "mongodbatlas_org_service_account" "this" {
  org_id                     = var.org_id
  name                       = "My organization service account"
  description                = "Service account in read-only"
  roles                      = ["ORG_READ_ONLY"]
  secret_expires_after_hours = 12
}