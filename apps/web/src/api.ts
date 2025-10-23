import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:9000/api/v1',
  withCredentials: true,   // 关键：携带 Cookie
  timeout: 8000,
})

export default api
