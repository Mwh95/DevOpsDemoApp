import React from 'react'
import ReactDOM from 'react-dom/client'
import { AuthProvider } from 'react-oidc-context'
import App from './App'
import './index.css'

const authority =
  import.meta.env.VITE_OIDC_AUTHORITY ??
  (typeof window !== 'undefined' ? window.location.origin + '/login/realms/master' : undefined)
const clientId = import.meta.env.VITE_OIDC_CLIENT_ID ?? 'map-app'
const redirectUri =
  import.meta.env.VITE_OIDC_REDIRECT_URI ?? (typeof window !== 'undefined' ? window.location.origin + '/' : '')

if (!authority) {
  throw new Error('Missing required VITE_OIDC_AUTHORITY configuration')
}

const oidcConfig = {
  authority,
  client_id: clientId,
  redirect_uri: redirectUri,
  scope: 'openid profile',
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AuthProvider {...oidcConfig}>
      <App />
    </AuthProvider>
  </React.StrictMode>
)
