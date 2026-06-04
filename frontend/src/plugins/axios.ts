import axios from 'axios'

// Set the base URL from environment variable or use relative path '/api' (inherits current domain/port automatically)
const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'
axios.defaults.baseURL = API_BASE_URL

// Add a request interceptor to include token in headers
axios.interceptors.request.use(
  function (config) {
    // Get token from localStorage
    const token = localStorage.getItem('token')
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    return config
  },
  function (error) {
    return Promise.reject(error)
  }
)

// Add a response interceptor for handling responses and errors
axios.interceptors.response.use(
  function (response) {
    return response
  },
  function (error) {
    // Handle specific error responses (e.g., unauthorized)
    if (error.response?.status === 401) {
      // Remove token and redirect to login if unauthenticated
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    
    return Promise.reject(error)
  }
)

export default axios