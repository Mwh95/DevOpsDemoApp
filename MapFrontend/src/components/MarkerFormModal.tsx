import { useState } from 'react'

interface MarkerFormModalProps {
  title: string
  initialLabel: string
  initialNote: string
  onSave: (label: string, note: string) => void
  onCancel: () => void
}

export function MarkerFormModal({
  title,
  initialLabel,
  initialNote,
  onSave,
  onCancel,
}: MarkerFormModalProps) {
  const [label, setLabel] = useState(initialLabel)
  const [note, setNote] = useState(initialNote)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave(label, note)
  }

  return (
    <div className="marker-form-modal" role="dialog" aria-modal="true">
      <div className="marker-form-inner">
        <h3>{title}</h3>
        <form onSubmit={handleSubmit}>
          <label htmlFor="marker-label">Label</label>
          <input
            id="marker-label"
            type="text"
            value={label}
            onChange={(e) => setLabel(e.target.value)}
            placeholder="e.g. Home, Office"
          />
          <label htmlFor="marker-note">Note</label>
          <textarea
            id="marker-note"
            value={note}
            onChange={(e) => setNote(e.target.value)}
            placeholder="Personal note..."
          />
          <div className="marker-form-actions">
            <button type="button" onClick={onCancel}>
              Cancel
            </button>
            <button type="submit">Save</button>
          </div>
        </form>
      </div>
    </div>
  )
}
