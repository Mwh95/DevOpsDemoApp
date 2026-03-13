import { useAuth } from 'react-oidc-context'
import { MapView } from './components/MapView'
import './App.css'

function App() {
  const auth = useAuth()

  if (auth.isLoading) {
    return (
      <div className="app-loading">
        <p>Loading...</p>
      </div>
    )
  }

  if (!auth.isAuthenticated || !auth.user) {
    return (
      <div className="app-auth">
        <p>You need to sign in to use the map.</p>
        <button type="button" onClick={() => auth.signinRedirect()}>
          Sign in
        </button>
      </div>
    )
  }

  const accessToken = (auth.user as { access_token?: string }).access_token

  return (
    <div className="app">
      <header className="app-header">
        <span>Map Markers</span>
        <button type="button" onClick={() => auth.signoutRedirect()}>
          Sign out
        </button>
      </header>
      <MapView accessToken={accessToken} />
    </div>
  )
}

export default App
