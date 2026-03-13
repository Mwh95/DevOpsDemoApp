interface EnvConfig {
  VITE_API_BASE: string
  VITE_OIDC_AUTHORITY: string
  VITE_OIDC_CLIENT_ID: string
  VITE_OIDC_REDIRECT_URI: string
}

declare global {
  interface Window {
    __ENV__?: Partial<EnvConfig>
  }
}

function get(key: keyof EnvConfig, fallback?: string): string {
  return window.__ENV__?.[key] ?? import.meta.env[key] ?? fallback ?? ''
}

export const config: EnvConfig = {
  VITE_API_BASE: get('VITE_API_BASE'),
  VITE_OIDC_AUTHORITY:
    get('VITE_OIDC_AUTHORITY') || window.location.origin + '/login/realms/master',
  VITE_OIDC_CLIENT_ID: get('VITE_OIDC_CLIENT_ID', 'map-app'),
  VITE_OIDC_REDIRECT_URI:
    get('VITE_OIDC_REDIRECT_URI') || window.location.origin + '/',
}
