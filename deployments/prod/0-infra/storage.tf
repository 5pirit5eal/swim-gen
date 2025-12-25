# Storage bucket for the exported pdfs from the backend
resource "google_storage_bucket" "exported_pdfs" {
  name     = "${var.project_id}-swim-gen-exported-pdfs"
  location = var.region
  project  = var.project_id

  lifecycle {
    prevent_destroy = true
  }

  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "Delete"
    }
  }

  depends_on = [google_project_service.apis]
}

# Storage bucket for public images
resource "google_storage_bucket" "public_images" {
  name     = "${var.project_id}-swim-gen-public-images"
  location = var.region
  project  = var.project_id

  lifecycle {
    prevent_destroy = true
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

# Upload images to the public bucket
resource "google_storage_bucket_object" "public_images_objects" {
  for_each = fileset("${path.module}/../../../data/images", "*.png")

  name   = each.value
  source = "${path.module}/../../../data/images/${each.value}"
  bucket = google_storage_bucket.public_images.name
}
