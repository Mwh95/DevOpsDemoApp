import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { MarkerLayer } from './MarkerLayer'

// Mock react-leaflet so we don't need a map context
vi.mock('react-leaflet', () => ({
  Marker: ({ children }: { children: React.ReactNode }) => <div data-testid="marker">{children}</div>,
  Popup: ({ children }: { children: React.ReactNode }) => <div data-testid="popup">{children}</div>,
}))

describe('MarkerLayer', () => {
  const markers = [
    {
      id: '1',
      user_id: 'u1',
      latitude: 52.52,
      longitude: 13.405,
      label: 'Home',
      note: 'My home',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    },
    {
      id: '2',
      user_id: 'u1',
      latitude: 52.53,
      longitude: 13.41,
      label: '',
      note: 'No label',
      created_at: '2024-01-02T00:00:00Z',
      updated_at: '2024-01-02T00:00:00Z',
    },
  ]

  it('renders a marker for each item', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()
    render(
      <MarkerLayer markers={markers} onEdit={onEdit} onDelete={onDelete} />
    )
    const markerEls = screen.getAllByTestId('marker')
    expect(markerEls).toHaveLength(2)
  })

  it('shows label and note in popup', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()
    render(
      <MarkerLayer markers={markers} onEdit={onEdit} onDelete={onDelete} />
    )
    const homes = screen.getAllByText('Home')
    expect(homes.length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('My home').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('(No label)').length).toBeGreaterThanOrEqual(1)
    expect(screen.getAllByText('No label').length).toBeGreaterThanOrEqual(1)
  })

  it('renders Edit and Delete buttons', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()
    render(
      <MarkerLayer markers={[markers[0]]} onEdit={onEdit} onDelete={onDelete} />
    )
    const editBtns = screen.getAllByRole('button', { name: /edit/i })
    const deleteBtns = screen.getAllByRole('button', { name: /delete/i })
    expect(editBtns.length).toBeGreaterThanOrEqual(1)
    expect(deleteBtns.length).toBeGreaterThanOrEqual(1)
  })

  it('renders no markers when markers is empty', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()
    render(
      <MarkerLayer markers={[]} onEdit={onEdit} onDelete={onDelete} />
    )
    const markerEls = screen.queryAllByTestId('marker')
    expect(markerEls).toHaveLength(0)
  })

  it('does not render note div when note is empty', () => {
    const onEdit = vi.fn()
    const onDelete = vi.fn()
    const markersNoNote = [
      { ...markers[0], id: 'x', label: 'NoNote', note: '' },
    ]
    const { container } = render(
      <MarkerLayer markers={markersNoNote} onEdit={onEdit} onDelete={onDelete} />
    )
    expect(screen.getByText('NoNote')).toBeInTheDocument()
    const popupNotes = container.querySelectorAll('.popup-note')
    expect(popupNotes).toHaveLength(0)
  })
})
