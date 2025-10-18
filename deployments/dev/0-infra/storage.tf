# Storage bucket for the exported pdfs from the backend
resource "google_storage_bucket" "exported_pdfs" {
  name          = "${var.project_id}-swim-gen-exported-pdfs"
  location      = var.region
  project       = var.project_id
  force_destroy = true

  lifecycle {
    prevent_destroy = false
  }

  depends_on = [google_project_service.apis]
}
