data "mongodbatlas_org_service_account" "this" {
  org_id    = var.org_id
  client_id = var.org_sa_id
}