import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Handle errors
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error.message)
  }
)

export const authService = {
  login: async (email, password) => {
    return api.post('/login', { email, password })
  },
  register: async (email, username, password) => {
    return api.post('/register', { email, username, password })
  },
  getProfile: async () => {
    return api.get('/user/profile')
  },
}

export const tradeService = {
  buy: async (asset_symbol, quantity, price) => {
    return api.post('/trade/buy', { asset_symbol, quantity, price })
  },
  sell: async (asset_symbol, quantity, price) => {
    return api.post('/trade/sell', { asset_symbol, quantity, price })
  },
  getHistory: async () => {
    return api.get('/trades/history')
  },
}

export const portfolioService = {
  getPortfolio: async () => {
    return api.get('/portfolio')
  },
  getHoldings: async () => {
    return api.get('/portfolio/holdings')
  },
}

export default api
