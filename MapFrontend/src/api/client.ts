import type { Marker, CreateMarkerInput, UpdateMarkerInput } from '../types/marker'

const API_BASE = import.meta.env.VITE_API_BASE ?? ''

function authHeaders(token: string | undefined): HeadersInit {
  const h: HeadersInit = { 'Content-Type': 'application/json' }
  if (token) {
    ;(h as Record<string, string>)['Authorization'] = `Bearer ${token}`
  }
  return h
}

export async function fetchMarkers(accessToken: string | undefined): Promise<Marker[]> {
  const res = await fetch(`${API_BASE}/api/markers`, { headers: authHeaders(accessToken) })
  if (res.status === 401) throw new Error('unauthorized')
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function createMarker(
  accessToken: string | undefined,
  data: CreateMarkerInput
): Promise<Marker> {
  const res = await fetch(`${API_BASE}/api/markers`, {
    method: 'POST',
    headers: authHeaders(accessToken),
    body: JSON.stringify(data),
  })
  if (res.status === 401) throw new Error('unauthorized')
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function updateMarker(
  accessToken: string | undefined,
  id: string,
  data: UpdateMarkerInput
): Promise<Marker> {
  const res = await fetch(`${API_BASE}/api/markers/${id}`, {
    method: 'PUT',
    headers: authHeaders(accessToken),
    body: JSON.stringify(data),
  })
  if (res.status === 401) throw new Error('unauthorized')
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function deleteMarker(
  accessToken: string | undefined,
  id: string
): Promise<void> {
  const res = await fetch(`${API_BASE}/api/markers/${id}`, {
    method: 'DELETE',
    headers: authHeaders(accessToken),
  })
  if (res.status === 401) throw new Error('unauthorized')
  if (!res.ok) throw new Error(await res.text())
}

export type { Marker, CreateMarkerInput, UpdateMarkerInput }
