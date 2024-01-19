import { BASE_URL } from "./config"

const PATH_PREFIX = 'v1/users'

export const signup = async (username, password) => {
  const body = JSON.stringify({ username, password }) 
  const response = await fetch(`${BASE_URL}/${PATH_PREFIX}/signup`, { method: 'POST', body })
  const result = await response.json()
  return result
}

export const signin = async (username, password) => {
  const body = JSON.stringify({ username, password })
  const response = await fetch(`${BASE_URL}/${PATH_PREFIX}/signin`, { method: 'POST', body })
  const result = await response.json()
  return result
}

export const verify = async () => {
  const response = await fetch(`${BASE_URL}/${PATH_PREFIX}/verify`, { method: 'POST', credentials: "include" })
  const result = await response.json()
  return result
}