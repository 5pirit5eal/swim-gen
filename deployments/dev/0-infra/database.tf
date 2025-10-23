resource "postgresql_extension" "pgvector" {
  name         = "vector"
  schema       = "extensions"
  drop_cascade = true
}

resource "postgresql_extension" "pg_cron" {
  name         = "pg_cron"
  schema       = "pg_catalog"
  drop_cascade = true
}
