import axios from "axios"

export const API_URL = `${process.env.HOST}:${process.env.API_PORT}` || 'http://localhost:8000'
export const BASE_URL = `${API_URL}/api`
export const DEFAULT_REQUEST_TIMEOUT = 5000
export const client = axios.create({ baseURL: BASE_URL, timeout: DEFAULT_REQUEST_TIMEOUT })