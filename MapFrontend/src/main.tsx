import React from 'react'
import ReactDOM from 'react-dom/client'
import { AuthProvider } from 'react-oidc-context'
import App from './App'
import './index.css'
import { config } from './config'

const oidcConfig = {
  authority: config.VITE_OIDC_AUTHORITY,
  client_id: config.VITE_OIDC_CLIENT_ID,
  redirect_uri: config.VITE_OIDC_REDIRECT_URI,
  scope: 'openid profile',
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AuthProvider {...oidcConfig}>
      <App />
    </AuthProvider>
  </React.StrictMode>
)
