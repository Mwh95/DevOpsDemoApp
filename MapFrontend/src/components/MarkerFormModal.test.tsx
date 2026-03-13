import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { MarkerFormModal } from './MarkerFormModal'

describe('MarkerFormModal', () => {
  it('renders title and fields', () => {
    const onSave = vi.fn()
    const onCancel = vi.fn()
    render(
      <MarkerFormModal
        title="New marker"
        initialLabel=""
        initialNote=""
        onSave={onSave}
        onCancel={onCancel}
      />
    )
    expect(screen.getByRole('heading', { name: 'New marker' })).toBeInTheDocument()
    expect(screen.getByLabelText(/label/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/note/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /save/i })).toBeInTheDocument()
  })

  it('calls onSave with label and note on submit', () => {
    const onSave = vi.fn()
    const onCancel = vi.fn()
    render(
      <MarkerFormModal
        title="Edit"
        initialLabel="Old"
        initialNote="Note"
        onSave={onSave}
        onCancel={onCancel}
      />
    )
    const labelInput = screen.getByLabelText(/label/i)
    const noteInput = screen.getByLabelText(/note/i)
    fireEvent.change(labelInput, { target: { value: 'Updated' } })
    fireEvent.change(noteInput, { target: { value: 'Updated note' } })
    fireEvent.click(screen.getByRole('button', { name: /^save$/i }))
    expect(onSave).toHaveBeenCalledWith('Updated', 'Updated note')
    expect(onCancel).not.toHaveBeenCalled()
  })

  it('calls onCancel when Cancel is clicked', () => {
    const onSave = vi.fn()
    const onCancel = vi.fn()
    render(
      <MarkerFormModal
        title="New"
        initialLabel=""
        initialNote=""
        onSave={onSave}
        onCancel={onCancel}
      />
    )
    fireEvent.click(screen.getByRole('button', { name: /^cancel$/i }))
    expect(onCancel).toHaveBeenCalled()
    expect(onSave).not.toHaveBeenCalled()
  })
})
