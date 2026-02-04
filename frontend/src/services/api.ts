import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface Pack {
  size: number
  quantity: number
}

export interface PackSizesResponse {
  sizes: number[]
}

export interface CalculateResponse {
  packs: Pack[]
}

export const getPackSizes = async (): Promise<number[]> => {
  const response = await api.get<PackSizesResponse>('/pack-sizes')
  return response.data.sizes
}

export const updatePackSizes = async (sizes: number[]): Promise<void> => {
  await api.post('/pack-sizes', { sizes })
}

export const calculatePacks = async (items: number): Promise<Pack[]> => {
  const response = await api.post<CalculateResponse>('/calculate', { items })
  return response.data.packs
}
