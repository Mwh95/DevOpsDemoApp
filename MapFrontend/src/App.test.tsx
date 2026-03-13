import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { useAuth } from 'react-oidc-context'
import App from './App'

vi.mock('react-oidc-context', () => ({
  useAuth: vi.fn(),
  AuthProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}))

vi.mock('./components/MapView', () => ({
  MapView: () => <div data-testid="map-view">MapView</div>,
}))

describe('App', () => {
  it('shows loading when auth is loading', () => {
    vi.mocked(useAuth).mockReturnValue({
      isLoading: true,
      isAuthenticated: false,
      user: null,
      signinRedirect: vi.fn(),
      signoutRedirect: vi.fn(),
    } as unknown as ReturnType<typeof useAuth>)
    render(<App />)
    expect(screen.getByText(/loading/i)).toBeInTheDocument()
  })

  it('shows sign in when not authenticated', () => {
    vi.mocked(useAuth).mockReturnValue({
      isLoading: false,
      isAuthenticated: false,
      user: null,
      signinRedirect: vi.fn(),
      signoutRedirect: vi.fn(),
    } as unknown as ReturnType<typeof useAuth>)
    render(<App />)
    expect(screen.getByText(/you need to sign in/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
  })

  it('shows MapView when authenticated', () => {
    const signoutRedirect = vi.fn()
    vi.mocked(useAuth).mockReturnValue({
      isLoading: false,
      isAuthenticated: true,
      user: { access_token: 'token' },
      signinRedirect: vi.fn(),
      signoutRedirect,
    } as unknown as ReturnType<typeof useAuth>)
    render(<App />)
    expect(screen.getByTestId('map-view')).toBeInTheDocument()
    expect(screen.getByText('Map Markers')).toBeInTheDocument()
    fireEvent.click(screen.getByRole('button', { name: /sign out/i }))
    expect(signoutRedirect).toHaveBeenCalled()
  })

  it('calls signinRedirect when Sign in is clicked', () => {
    const signinRedirect = vi.fn()
    vi.mocked(useAuth).mockReturnValue({
      isLoading: false,
      isAuthenticated: false,
      user: null,
      signinRedirect,
      signoutRedirect: vi.fn(),
    } as unknown as ReturnType<typeof useAuth>)
    render(<App />)
    fireEvent.click(screen.getByRole('button', { name: /sign in/i }))
    expect(signinRedirect).toHaveBeenCalled()
  })
})
