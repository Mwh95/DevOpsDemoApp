variable "keycloak_url" {
  description = "Base URL of the Keycloak instance (before the REST base path)."
  type        = string
}

variable "keycloak_base_path" {
  description = "REST base path; matches KC_HTTP_RELATIVE_PATH (e.g. /login)."
  type        = string
  default     = "/login"
}

variable "keycloak_admin_user" {
  description = "Bootstrap admin username in the master realm."
  type        = string
  default     = "tmpadmin"
}

variable "keycloak_admin_password" {
  description = "Bootstrap admin password in the master realm."
  type        = string
}

variable "realm" {
  description = "Realm to create for the Map Markers app."
  type        = string
  default     = "users"
}

variable "client_id" {
  description = "OIDC client id used by the SPA."
  type        = string
  default     = "map-app"
}

variable "redirect_uris" {
  description = "Valid redirect URIs for the public SPA client."
  type        = list(string)
  default     = ["http://localhost:58080/*", "http://localhost:5173/*"]
}

variable "web_origins" {
  description = "Allowed CORS web origins for the SPA client."
  type        = list(string)
  default     = ["+"]
}

variable "test_password" {
  description = "Password assigned to the seeded test users."
  type        = string
  default     = "Test1234!"
}
