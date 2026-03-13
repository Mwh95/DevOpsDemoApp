import { Marker, Popup } from 'react-leaflet'
import L from 'leaflet'
import type { Marker as MarkerType } from '../types/marker'

// Fix default icon
const icon = L.icon({
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
  iconSize: [25, 41],
  iconAnchor: [12, 41],
})

interface MarkerLayerProps {
  markers: MarkerType[]
  onEdit: (m: MarkerType) => void
  onDelete: (id: string) => void
}

export function MarkerLayer({ markers, onEdit, onDelete }: MarkerLayerProps) {
  return (
    <>
      {markers.map((m) => (
        <Marker
          key={m.id}
          position={[m.latitude, m.longitude]}
          icon={icon}
        >
          <Popup>
            <div className="marker-popup">
              <div className="popup-label">{m.label || '(No label)'}</div>
              {m.note && <div className="popup-note">{m.note}</div>}
              <div className="marker-popup-actions">
                <button type="button" onClick={() => onEdit(m)}>
                  Edit
                </button>
                <button type="button" onClick={() => onDelete(m.id)}>
                  Delete
                </button>
              </div>
            </div>
          </Popup>
        </Marker>
      ))}
    </>
  )
}
