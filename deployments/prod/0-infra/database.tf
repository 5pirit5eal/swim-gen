resource "postgresql_extension" "pgvector" {
  name         = "vector"
  schema       = "extensions"
  drop_cascade = true
}
