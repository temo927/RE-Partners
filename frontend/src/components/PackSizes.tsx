import { useState, useEffect } from 'react'
import { getPackSizes, updatePackSizes } from '../services/api'
import './PackSizes.css'

export default function PackSizes() {
  const [sizes, setSizes] = useState<string[]>(['', '', ''])
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState<string>('')

  useEffect(() => {
    loadPackSizes()
  }, [])

  const loadPackSizes = async () => {
    try {
      const packSizes = await getPackSizes()
      const newSizes = [...sizes]
      packSizes.forEach((size, index) => {
        if (index < 3) {
          newSizes[index] = size.toString()
        }
      })
      setSizes(newSizes)
    } catch (error) {
      console.error('Failed to load pack sizes:', error)
    }
  }

  const handleSizeChange = (index: number, value: string) => {
    const newSizes = [...sizes]
    newSizes[index] = value
    setSizes(newSizes)
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
              <input
                key={index}
                type="number"
                className="pack-size-input"
                placeholder={`Pack size ${index + 1}`}
                value={size}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => handleSizeChange(index, e.target.value)}
                min="1"
              />
            ))}
          </div>
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
