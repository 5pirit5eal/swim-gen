########################################
# Database Extensions
########################################

resource "postgresql_extension" "pgvector_public" {
  name         = "vector"
  schema       = "public"
  drop_cascade = true
}

resource "postgresql_extension" "pgvector_extensions" {
  name         = "vector"
  schema       = "extensions"
  drop_cascade = true
}

########################################
# User Creation
########################################

resource "postgresql_role" "backend_user" {
  name                      = var.dbusers.backend
  login                     = true
  password                  = data.google_secret_manager_secret_version_access.dbpassword_user.secret_data
  create_database           = false
  create_role               = false
  bypass_row_level_security = true
  valid_until               = "infinity"
}

resource "postgresql_role" "frontend_user" {
  name                      = var.dbusers.frontend
  login                     = true
  password                  = data.google_secret_manager_secret_version_access.dbpassword_user.secret_data
  create_database           = false
  create_role               = false
  bypass_row_level_security = true
  valid_until               = "infinity"
}

# ########################################
# # Schemas Creation & grant creation inside schema
# ########################################

resource "postgresql_schema" "schema" {
  database      = "postgres"
  name          = "public"
  owner         = "postgres"
  if_not_exists = true
  drop_cascade  = true

  depends_on = [
    postgresql_role.backend_user,
    postgresql_role.frontend_user,
  ]
}

resource "postgresql_grant" "grant_roles_schema" {

  for_each = { "postgres" = ["USAGE", "CREATE"], "${var.dbusers.backend}" = ["USAGE", "CREATE"], "${var.dbusers.frontend}" = ["USAGE"] }

  database    = "postgres"
  schema      = postgresql_schema.schema.name
  role        = each.key
  object_type = "schema"
  privileges  = try(each.value, null)
}

########################################
# Creation of grants for each role
########################################
resource "postgresql_grant" "backend_privileges" {
  for_each = {
    "table"    = ["SELECT", "INSERT", "UPDATE", "DELETE", "TRUNCATE", "REFERENCES", "TRIGGER"],
    "database" = ["CONNECT"],
    "sequence" = ["USAGE", "SELECT"],
  }
  database    = "postgres"
  schema      = postgresql_schema.schema.name
  role        = var.dbusers.backend
  object_type = each.key
  privileges  = each.value

  depends_on = [
    postgresql_grant.grant_roles_schema,
    postgresql_role.backend_user,
  ]
}

resource "postgresql_grant" "frontend_privileges" {
  for_each = {
    "table"    = ["SELECT", "REFERENCES", "TRIGGER"],
    "database" = ["CONNECT"],
    "sequence" = ["USAGE", "SELECT"],
  }
  database    = "postgres"
  schema      = postgresql_schema.schema.name
  role        = var.dbusers.frontend
  object_type = each.key
  privileges  = each.value

  depends_on = [
    postgresql_grant.grant_roles_schema,
    postgresql_role.backend_user,
  ]
}

########################################
# REVOKE CREATE ON SCHEMA public FROM PUBLIC;
# Because by default, the default privileges allow any user ("public")
# to create table inside "public" schema
########################################
resource "postgresql_grant" "revoke_create_public" {
  database    = "postgres"
  schema      = postgresql_schema.schema.name
  role        = "public"
  object_type = "schema"
  privileges  = []

  depends_on = [
    postgresql_grant.backend_privileges,
    postgresql_grant.frontend_privileges,
  ]
}
