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

# Storage bucket for public images
resource "google_storage_bucket" "public_images" {
  name          = "${var.project_id}-swim-gen-public-images"
  location      = var.region
  project       = var.project_id
  force_destroy = true

  lifecycle {
    prevent_destroy = false
  }

  uniform_bucket_level_access = true

  cors {
    origin          = ["*"]
    method          = ["GET", "HEAD", "OPTIONS"]
    response_header = ["*"]
    max_age_seconds = 3600
  }

  depends_on = [google_project_service.apis]
}

# Make the bucket public
resource "google_storage_bucket_iam_member" "public_images_iam" {
  bucket = google_storage_bucket.public_images.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}
