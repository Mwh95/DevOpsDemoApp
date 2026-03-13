import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  createMarker,
  deleteMarker,
  fetchMarkers,
  updateMarker,
} from './client'

describe('api client', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('fetchMarkers returns list with token', async () => {
    const token = 'access-token'
    const mockMarkers = [
      {
        id: '1',
        user_id: 'u1',
        latitude: 52.52,
        longitude: 13.405,
        label: 'Home',
        note: 'Note',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
    ]
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => mockMarkers,
    })
    const result = await fetchMarkers(token)
    expect(result).toEqual(mockMarkers)
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/markers'),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: 'Bearer access-token',
        }),
      })
    )
  })

  it('fetchMarkers throws on 401', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 401,
      text: async () => 'unauthorized',
    })
    await expect(fetchMarkers('token')).rejects.toThrow('unauthorized')
  })

  it('createMarker sends POST with body', async () => {
    const created = {
      id: 'new-id',
      user_id: 'u1',
      latitude: 52.52,
      longitude: 13.405,
      label: 'Office',
      note: 'Work',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 201,
      json: async () => created,
    })
    const result = await createMarker('token', {
      latitude: 52.52,
      longitude: 13.405,
      label: 'Office',
      note: 'Work',
    })
    expect(result).toEqual(created)
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/markers'),
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({
          latitude: 52.52,
          longitude: 13.405,
          label: 'Office',
          note: 'Work',
        }),
      })
    )
  })

  it('updateMarker sends PUT with id', async () => {
    const updated = {
      id: 'id1',
      user_id: 'u1',
      latitude: 52.52,
      longitude: 13.405,
      label: 'New Label',
      note: 'New note',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-02T00:00:00Z',
    }
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => updated,
    })
    const result = await updateMarker('token', 'id1', {
      label: 'New Label',
      note: 'New note',
    })
    expect(result).toEqual(updated)
    expect(fetch).toHaveBeenCalledWith(
      expect.stringMatching(/\/api\/markers\/id1$/),
      expect.objectContaining({
        method: 'PUT',
        body: JSON.stringify({ label: 'New Label', note: 'New note' }),
      })
    )
  })

  it('deleteMarker sends DELETE', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 204,
      text: async () => '',
    })
    await deleteMarker('token', 'id1')
    expect(fetch).toHaveBeenCalledWith(
      expect.stringMatching(/\/api\/markers\/id1$/),
      expect.objectContaining({ method: 'DELETE' })
    )
  })

  it('createMarker throws on 401', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 401,
      text: async () => 'unauthorized',
    })
    await expect(
      createMarker('token', { latitude: 1, longitude: 1, label: '', note: '' })
    ).rejects.toThrow('unauthorized')
  })

  it('createMarker throws when not ok', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 400,
      text: async () => 'bad request',
    })
    await expect(
      createMarker('token', { latitude: 1, longitude: 1, label: '', note: '' })
    ).rejects.toThrow('bad request')
  })

  it('updateMarker throws when not ok', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 500,
      text: async () => 'server error',
    })
    await expect(
      updateMarker('token', 'id1', { label: 'x' })
    ).rejects.toThrow('server error')
  })

  it('updateMarker throws on 401', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 401,
      text: async () => 'unauthorized',
    })
    await expect(
      updateMarker('token', 'id1', { label: 'x' })
    ).rejects.toThrow('unauthorized')
  })

  it('deleteMarker throws on 401', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 401,
      text: async () => 'unauthorized',
    })
    await expect(deleteMarker('token', 'id1')).rejects.toThrow('unauthorized')
  })

  it('fetchMarkers without token does not set Authorization header', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: true,
      status: 200,
      json: async () => [],
    })
    await fetchMarkers(undefined)
    const call = (fetch as ReturnType<typeof vi.fn>).mock.calls[0]
    expect(call[1].headers).toEqual({ 'Content-Type': 'application/json' })
  })

  it('fetchMarkers throws when not ok and status is not 401', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 500,
      text: async () => 'server error',
    })
    await expect(fetchMarkers('token')).rejects.toThrow('server error')
  })

  it('deleteMarker throws when not ok and status is not 401', async () => {
    ;(fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
      ok: false,
      status: 500,
      text: async () => 'server error',
    })
    await expect(deleteMarker('token', 'id1')).rejects.toThrow('server error')
  })
})
