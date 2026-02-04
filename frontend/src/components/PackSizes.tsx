import { useState, useEffect } from 'react'
import { getPackSizes, updatePackSizes } from '../services/api'
import './PackSizes.css'

export default function PackSizes() {
  const [sizes, setSizes] = useState<string[]>(['23', '31', '53'])
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState<string>('')

  useEffect(() => {
    loadPackSizes()
  }, [])

  const loadPackSizes = async () => {
    try {
      const packSizes = await getPackSizes()
      if (packSizes.length > 0) {
        setSizes(packSizes.map(size => size.toString()))
      } else {
        setSizes(['23', '31', '53'])
      }
    } catch (error) {
      console.error('Failed to load pack sizes:', error)
      setSizes(['23', '31', '53'])
    }
  }

  const handleSizeChange = (index: number, value: string) => {
    const newSizes = [...sizes]
    newSizes[index] = value
    setSizes(newSizes)
  }

  const handleAddSize = () => {
    setSizes([...sizes, ''])
  }

  const handleRemoveSize = (index: number) => {
    if (sizes.length > 1) {
      const newSizes = sizes.filter((_, i) => i !== index)
      setSizes(newSizes)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setMessage('')

    try {
      const packSizesArray = sizes
        .map((size: string) => parseInt(size.trim()))
        .filter((size: number) => !isNaN(size) && size > 0)

      if (packSizesArray.length === 0) {
        setMessage('Please enter at least one valid pack size')
        setLoading(false)
        return
      }

      await updatePackSizes(packSizesArray)
      setMessage('Pack sizes updated successfully')
      setTimeout(() => setMessage(''), 3000)
    } catch (error: any) {
      setMessage(error.response?.data?.error || 'Failed to update pack sizes')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="pack-sizes-container">
      <div className="pack-sizes-box">
        <h2 className="pack-sizes-header">Pack Sizes</h2>
        <form onSubmit={handleSubmit}>
          <div className="pack-sizes-inputs">
            {sizes.map((size: string, index: number) => (
              <div key={index} className="pack-size-input-wrapper">
                <input
                  type="number"
                  className="pack-size-input"
                  placeholder={`Pack size ${index + 1}`}
                  value={size}
                  onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleSizeChange(index, e.target.value)}
                  min="1"
                />
                {sizes.length > 1 && (
                  <button
                    type="button"
                    className="remove-button"
                    onClick={() => handleRemoveSize(index)}
                    aria-label="Remove pack size"
                  >
                    ×
                  </button>
                )}
              </div>
            ))}
          </div>
          <button
            type="button"
            className="add-button"
            onClick={handleAddSize}
          >
            + Add Pack Size
          </button>
          {message && (
            <div className={`message ${message.includes('success') ? 'success' : 'error'}`}>
              {message}
            </div>
          )}
          <button type="submit" className="submit-button" disabled={loading}>
            {loading ? 'Submitting...' : 'Submit pack sizes change'}
          </button>
        </form>
      </div>
    </div>
  )
}
