terraform {
  required_version = ">= 1.6"
  required_providers {
    keycloak = {
      source  = "keycloak/keycloak"
      version = ">= 5.8.0"
    }
  }
}

# Authenticates against the master realm using the admin-cli client and the
# bootstrap admin credentials (password grant).
provider "keycloak" {
  client_id = "admin-cli"
  username  = var.keycloak_admin_user
  password  = var.keycloak_admin_password
  url       = var.keycloak_url
  base_path = var.keycloak_base_path
}

resource "keycloak_realm" "users" {
  realm        = var.realm
  enabled      = true
  ssl_required = "none"
}

# Public SPA client. direct_access_grants_enabled is intentionally on so the
# system tests can obtain tokens via the password grant for API checks; the
# browser uses the standard (authorization code + PKCE) flow.
resource "keycloak_openid_client" "map_app" {
  realm_id                       = keycloak_realm.users.id
  client_id                      = var.client_id
  name                           = "Map Markers SPA"
  enabled                        = true
  access_type                    = "PUBLIC"
  standard_flow_enabled          = true
  implicit_flow_enabled          = false
  direct_access_grants_enabled   = true
  valid_redirect_uris            = var.redirect_uris
  valid_post_logout_redirect_uris = ["+"]
  web_origins                    = var.web_origins
  pkce_code_challenge_method     = "S256"
}

resource "keycloak_user" "testuser" {
  realm_id       = keycloak_realm.users.id
  username       = "testuser"
  enabled        = true
  email          = "testuser@example.com"
  email_verified = true
  first_name     = "Test"
  last_name      = "User"

  initial_password {
    value     = var.test_password
    temporary = false
  }
}

resource "keycloak_user" "otheruser" {
  realm_id       = keycloak_realm.users.id
  username       = "otheruser"
  enabled        = true
  email          = "otheruser@example.com"
  email_verified = true
  first_name     = "Other"
  last_name      = "User"

  initial_password {
    value     = var.test_password
    temporary = false
  }
}
