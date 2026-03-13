import { useCallback, useEffect, useRef, useState } from 'react'
import { MapContainer, TileLayer, useMapEvents } from 'react-leaflet'
import L from 'leaflet'
import type { Marker as MarkerType } from '../types/marker'
import * as api from '../api/client'
import { MarkerFormModal } from './MarkerFormModal'
import { MarkerLayer } from './MarkerLayer'
import 'leaflet/dist/leaflet.css'

// Fix default icon in Leaflet with bundlers
delete (L.Icon.Default.prototype as unknown as { _getIconUrl?: unknown })._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
})

interface MapViewProps {
  accessToken: string | undefined
}

export function MapView({ accessToken }: MapViewProps) {
  const [markers, setMarkers] = useState<MarkerType[]>([])
  const [editMode, setEditMode] = useState(false)
  const [pendingMarker, setPendingMarker] = useState<{ lat: number; lng: number } | null>(null)
  const [editingMarker, setEditingMarker] = useState<MarkerType | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const accessTokenRef = useRef(accessToken)
  accessTokenRef.current = accessToken

  const loadMarkers = useCallback(async () => {
    if (!accessToken) return
    setLoading(true)
    setError(null)
    try {
      const list = await api.fetchMarkers(accessToken)
      setMarkers(list)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load markers')
    } finally {
      setLoading(false)
    }
  }, [accessToken])

  useEffect(() => {
    loadMarkers()
  }, [loadMarkers])

  const handleMapClick = useCallback(
    (lat: number, lng: number) => {
      if (!editMode || !accessToken) return
      setPendingMarker({ lat, lng })
    },
    [editMode, accessToken]
  )

  const handleSaveNew = useCallback(
    async (label: string, note: string) => {
      const token = accessTokenRef.current
      if (!pendingMarker || !token) return
      try {
        const created = await api.createMarker(token, {
          latitude: pendingMarker.lat,
          longitude: pendingMarker.lng,
          label,
          note,
        })
        setMarkers((prev) => [created, ...prev])
        setPendingMarker(null)
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Failed to create marker')
      }
    },
    [pendingMarker]
  )

  const handleUpdate = useCallback(
    async (id: string, label: string, note: string) => {
      const token = accessTokenRef.current
      if (!token) return
      try {
        const updated = await api.updateMarker(token, id, { label, note })
        setMarkers((prev) => prev.map((m) => (m.id === id ? updated : m)))
        setEditingMarker(null)
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Failed to update marker')
      }
    },
    []
  )

  const handleDelete = useCallback(
    async (id: string) => {
      const token = accessTokenRef.current
      if (!token) return
      try {
        await api.deleteMarker(token, id)
        setMarkers((prev) => prev.filter((m) => m.id !== id))
        setEditingMarker(null)
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Failed to delete marker')
      }
    },
    []
  )

  return (
    <div className="map-container">
      <div className="map-toolbar">
        <label>
          <input
            type="checkbox"
            checked={editMode}
            onChange={(e) => setEditMode(e.target.checked)}
          />
          Edit mode (click map to add marker)
        </label>
        {editMode && (
          <button
            type="button"
            onClick={() => setPendingMarker({ lat: 52.52, lng: 13.405 })}
            aria-label="Add marker at map center"
          >
            Add marker at center
          </button>
        )}
      </div>
      {error && (
        <div style={{ position: 'absolute', top: 50, left: 10, zIndex: 1000, background: '#fee', padding: '8px 12px', borderRadius: 4 }}>
          {error}
        </div>
      )}
      {loading && (
        <div style={{ position: 'absolute', top: 50, right: 10, zIndex: 1000, background: '#fff', padding: '4px 8px', borderRadius: 4 }}>
          Loading…
        </div>
      )}
      <MapContainer
        center={[52.52, 13.405]}
        zoom={10}
        style={{ height: '100%', width: '100%' }}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <MapClickHandler onMapClick={handleMapClick} />
        <MarkerLayer
          markers={markers}
          onEdit={setEditingMarker}
          onDelete={handleDelete}
        />
      </MapContainer>

      {pendingMarker && (
        <MarkerFormModal
          title="New marker"
          initialLabel=""
          initialNote=""
          onSave={handleSaveNew}
          onCancel={() => setPendingMarker(null)}
        />
      )}
      {editingMarker && (
        <MarkerFormModal
          title="Edit marker"
          initialLabel={editingMarker.label}
          initialNote={editingMarker.note}
          onSave={(label, note) => handleUpdate(editingMarker.id, label, note)}
          onCancel={() => setEditingMarker(null)}
        />
      )}
    </div>
  )
}

function MapClickHandler({
  onMapClick,
}: {
  onMapClick: (lat: number, lng: number) => void
}) {
  useMapEvents({
    click(e) {
      onMapClick(e.latlng.lat, e.latlng.lng)
    },
  })
  return null
}
