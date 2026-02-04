import { useState } from 'react'
import { calculatePacks } from '../services/api'
import './Calculator.css'

interface Pack {
  size: number
  quantity: number
}

export default function Calculator() {
  const [items, setItems] = useState<string>('500000')
  const [packs, setPacks] = useState<Pack[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string>('')

  const handleCalculate = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setPacks([])

    try {
      const itemsNum = parseInt(items.trim())
      if (isNaN(itemsNum) || itemsNum <= 0) {
        setError('Please enter a valid number of items')
        setLoading(false)
        return
      }

      const result = await calculatePacks(itemsNum)
      setPacks(result)
      setError('')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to calculate packs')
      setPacks([])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="calculator-container">
      <h2 className="calculator-header">Calculate packs for order</h2>
      <form onSubmit={handleCalculate}>
        <div className="items-input-group">
          <label htmlFor="items" className="items-label">Items:</label>
          <input
            id="items"
            type="number"
            className="items-input"
            value={items}
            onChange={(e) => setItems(e.target.value)}
            min="1"
          />
          <button type="submit" className="calculate-button" disabled={loading}>
            {loading ? 'Calculating...' : 'Calculate'}
          </button>
        </div>
        {error && <div className="error-message">{error}</div>}
      </form>

      {packs.length > 0 && (
        <div className="results-table">
          <div className="table-header">
            <div className="table-cell header-cell">Pack</div>
            <div className="table-cell header-cell">Quantity</div>
          </div>
          {packs.map((pack, index) => (
            <div key={index} className="table-row">
              <div className="table-cell">
                <input
                  type="text"
                  className="result-input"
                  value={pack.size}
                  readOnly
                />
              </div>
              <div className="table-cell">
                <input
                  type="text"
                  className="result-input"
                  value={pack.quantity}
                  readOnly
                />
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
