export interface Marker {
  id: string
  user_id: string
  latitude: number
  longitude: number
  label: string
  note: string
  created_at: string
  updated_at: string
}

export interface CreateMarkerInput {
  latitude: number
  longitude: number
  label: string
  note: string
}

export interface UpdateMarkerInput {
  label?: string
  note?: string
}
