import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { MapView } from './MapView'
import * as api from '../api/client'

vi.mock('../api/client')

const mockMapClickHandlers: { click?: (e: { latlng: { lat: number; lng: number } }) => void } = {}
vi.mock('react-leaflet', () => ({
  MapContainer: ({ children }: { children: React.ReactNode }) => <div data-testid="map-container">{children}</div>,
  TileLayer: () => null,
  useMapEvents: (handlers: { click?: (e: { latlng: { lat: number; lng: number } }) => void }) => {
    mockMapClickHandlers.click = handlers.click
    return null
  },
  Marker: ({ children }: { children: React.ReactNode }) => <div data-testid="marker">{children}</div>,
  Popup: ({ children }: { children: React.ReactNode }) => <div data-testid="popup">{children}</div>,
}))

const mockMarker = {
  id: 'm1',
  user_id: 'u1',
  latitude: 52.52,
  longitude: 13.405,
  label: 'Home',
  note: 'My note',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
}

describe('MapView', () => {
  it('renders map and edit mode toggle', () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    render(<MapView accessToken="token" />)
    expect(screen.getByTestId('map-container')).toBeInTheDocument()
    expect(screen.getByRole('checkbox', { name: /edit mode/i })).toBeInTheDocument()
  })

  it('renders with no token', () => {
    render(<MapView accessToken={undefined} />)
    expect(screen.getByTestId('map-container')).toBeInTheDocument()
  })

  it('calls fetchMarkers on mount when token present', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(api.fetchMarkers).toHaveBeenCalledWith('token')
    })
  })

  it('shows loading then markers list', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    render(<MapView accessToken="token" />)
    expect(screen.getByText(/loading/i)).toBeInTheDocument()
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    expect(screen.queryByText(/loading/i)).not.toBeInTheDocument()
  })

  it('shows error when fetchMarkers fails', async () => {
    vi.mocked(api.fetchMarkers).mockRejectedValue(new Error('network error'))
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('network error')).toBeInTheDocument()
    })
  })

  it('shows generic error when fetch fails with non-Error', async () => {
    vi.mocked(api.fetchMarkers).mockRejectedValue('string error')
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Failed to load markers')).toBeInTheDocument()
    })
  })

  it('toggles edit mode and shows add-marker button', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(api.fetchMarkers).toHaveBeenCalled()
    })
    expect(screen.queryByRole('button', { name: /add marker at map center/i })).not.toBeInTheDocument()
    fireEvent.click(screen.getByRole('checkbox', { name: /edit mode/i }))
    expect(screen.getByRole('button', { name: /add marker at map center/i })).toBeInTheDocument()
  })

  it('create flow: add marker at center opens modal, save calls createMarker', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    vi.mocked(api.createMarker).mockResolvedValue({ ...mockMarker, id: 'new-id', label: 'Office', note: 'Work' })
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(api.fetchMarkers).toHaveBeenCalled()
    })
    fireEvent.click(screen.getByRole('checkbox', { name: /edit mode/i }))
    fireEvent.click(screen.getByRole('button', { name: /add marker at map center/i }))
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /new marker/i })).toBeInTheDocument()
    })
    fireEvent.change(screen.getByLabelText(/label/i), { target: { value: 'Office' } })
    fireEvent.change(screen.getByLabelText(/note/i), { target: { value: 'Work' } })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    await waitFor(() => {
      expect(api.createMarker).toHaveBeenCalledWith('token', expect.objectContaining({ label: 'Office', note: 'Work', latitude: 52.52, longitude: 13.405 }))
    })
  })

  it('edit flow: click Edit on marker opens modal, save calls updateMarker', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    vi.mocked(api.updateMarker).mockResolvedValue({ ...mockMarker, label: 'Updated', note: 'Updated note' })
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    fireEvent.click(screen.getAllByRole('button', { name: /edit/i })[0])
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /edit marker/i })).toBeInTheDocument()
    })
    fireEvent.change(screen.getByLabelText(/label/i), { target: { value: 'Updated' } })
    fireEvent.change(screen.getByLabelText(/note/i), { target: { value: 'Updated note' } })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    await waitFor(() => {
      expect(api.updateMarker).toHaveBeenCalledWith('token', 'm1', { label: 'Updated', note: 'Updated note' })
    })
  })

  it('edit one of two markers updates only that marker in list', async () => {
    const marker2 = { ...mockMarker, id: 'm2', label: 'Other' }
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker, marker2])
    vi.mocked(api.updateMarker).mockResolvedValue({ ...mockMarker, label: 'Updated', note: 'New note' })
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
      expect(screen.getByText('Other')).toBeInTheDocument()
    })
    fireEvent.click(screen.getAllByRole('button', { name: /edit/i })[0])
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /edit marker/i })).toBeInTheDocument()
    })
    fireEvent.change(screen.getByLabelText(/label/i), { target: { value: 'Updated' } })
    fireEvent.change(screen.getByLabelText(/note/i), { target: { value: 'New note' } })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    await waitFor(() => {
      expect(api.updateMarker).toHaveBeenCalledWith('token', 'm1', { label: 'Updated', note: 'New note' })
    })
    expect(screen.getByText('Updated')).toBeInTheDocument()
    expect(screen.getByText('Other')).toBeInTheDocument()
  })

  it('delete flow: click Delete calls deleteMarker', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    vi.mocked(api.deleteMarker).mockResolvedValue()
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    fireEvent.click(screen.getAllByRole('button', { name: /delete/i })[0])
    await waitFor(() => {
      expect(api.deleteMarker).toHaveBeenCalledWith('token', 'm1')
    })
  })

  it('create marker error shows error message', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    vi.mocked(api.createMarker).mockRejectedValue(new Error('create failed'))
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(api.fetchMarkers).toHaveBeenCalled()
    })
    fireEvent.click(screen.getByRole('checkbox', { name: /edit mode/i }))
    fireEvent.click(screen.getByRole('button', { name: /add marker at map center/i }))
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /new marker/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    await waitFor(() => {
      expect(screen.getByText('create failed')).toBeInTheDocument()
    })
  })

  it('cancel pending marker closes modal', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(api.fetchMarkers).toHaveBeenCalled()
    })
    fireEvent.click(screen.getByRole('checkbox', { name: /edit mode/i }))
    fireEvent.click(screen.getByRole('button', { name: /add marker at map center/i }))
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /new marker/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /^cancel$/i }))
    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: /new marker/i })).not.toBeInTheDocument()
    })
  })

  it('update marker error shows error message', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    vi.mocked(api.updateMarker).mockRejectedValue(new Error('update failed'))
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    fireEvent.click(screen.getAllByRole('button', { name: /edit/i })[0])
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /edit marker/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    await waitFor(() => {
      expect(screen.getByText('update failed')).toBeInTheDocument()
    })
  })

  it('update marker non-Error rejection shows generic message', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    vi.mocked(api.updateMarker).mockRejectedValue('string error')
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    fireEvent.click(screen.getAllByRole('button', { name: /edit/i })[0])
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /edit marker/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    await waitFor(() => {
      expect(screen.getByText('Failed to update marker')).toBeInTheDocument()
    })
  })

  it('delete marker error shows error message', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    vi.mocked(api.deleteMarker).mockRejectedValue(new Error('delete failed'))
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    fireEvent.click(screen.getAllByRole('button', { name: /delete/i })[0])
    await waitFor(() => {
      expect(screen.getByText('delete failed')).toBeInTheDocument()
    })
  })

  it('map click in edit mode opens pending marker modal', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(api.fetchMarkers).toHaveBeenCalled()
    })
    fireEvent.click(screen.getByRole('checkbox', { name: /edit mode/i }))
    mockMapClickHandlers.click?.({ latlng: { lat: 52.52, lng: 13.405 } })
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /new marker/i })).toBeInTheDocument()
    })
  })

  it('map click when edit mode off does not open modal', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(api.fetchMarkers).toHaveBeenCalled()
    })
    expect(screen.queryByRole('checkbox', { name: /edit mode/i })).toBeInTheDocument()
    mockMapClickHandlers.click?.({ latlng: { lat: 52.52, lng: 13.405 } })
    await waitFor(() => {})
    expect(screen.queryByRole('heading', { name: /new marker/i })).not.toBeInTheDocument()
  })

  it('map click with no token does not open modal', async () => {
    render(<MapView accessToken={undefined} />)
    mockMapClickHandlers.click?.({ latlng: { lat: 52.52, lng: 13.405 } })
    expect(screen.queryByRole('heading', { name: /new marker/i })).not.toBeInTheDocument()
  })

  it('cancel edit marker modal closes modal', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    fireEvent.click(screen.getAllByRole('button', { name: /edit/i })[0])
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /edit marker/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /^cancel$/i }))
    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: /edit marker/i })).not.toBeInTheDocument()
    })
  })

  it('create marker non-Error rejection shows generic message', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    vi.mocked(api.createMarker).mockRejectedValue('create failed')
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(api.fetchMarkers).toHaveBeenCalled()
    })
    fireEvent.click(screen.getByRole('checkbox', { name: /edit mode/i }))
    fireEvent.click(screen.getByRole('button', { name: /add marker at map center/i }))
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /new marker/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    await waitFor(() => {
      expect(screen.getByText('Failed to create marker')).toBeInTheDocument()
    })
  })

  it('delete marker non-Error rejection shows generic message', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    vi.mocked(api.deleteMarker).mockRejectedValue('delete failed')
    render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    fireEvent.click(screen.getAllByRole('button', { name: /delete/i })[0])
    await waitFor(() => {
      expect(screen.getByText('Failed to delete marker')).toBeInTheDocument()
    })
  })

  it('handleSaveNew with no token does not call createMarker', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([])
    vi.mocked(api.createMarker).mockClear()
    const { rerender } = render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(api.fetchMarkers).toHaveBeenCalled()
    })
    fireEvent.click(screen.getByRole('checkbox', { name: /edit mode/i }))
    fireEvent.click(screen.getByRole('button', { name: /add marker at map center/i }))
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /new marker/i })).toBeInTheDocument()
    })
    rerender(<MapView accessToken={undefined} />)
    fireEvent.change(screen.getByLabelText(/label/i), { target: { value: 'X' } })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    await waitFor(() => {})
    expect(api.createMarker).not.toHaveBeenCalled()
  })

  it('handleUpdate with no token does not call updateMarker', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    vi.mocked(api.updateMarker).mockClear()
    const { rerender } = render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    rerender(<MapView accessToken={undefined} />)
    fireEvent.click(screen.getAllByRole('button', { name: /edit/i })[0])
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: /edit marker/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    await waitFor(() => {})
    expect(api.updateMarker).not.toHaveBeenCalled()
  })

  it('handleDelete with no token does not call deleteMarker', async () => {
    vi.mocked(api.fetchMarkers).mockResolvedValue([mockMarker])
    vi.mocked(api.deleteMarker).mockClear()
    const { rerender } = render(<MapView accessToken="token" />)
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
    rerender(<MapView accessToken={undefined} />)
    fireEvent.click(screen.getAllByRole('button', { name: /delete/i })[0])
    await waitFor(() => {})
    expect(api.deleteMarker).not.toHaveBeenCalled()
  })
})
