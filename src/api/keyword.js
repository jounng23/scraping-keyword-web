import { BASE_URL } from "./config"

const PATH_PREFIX = 'v1/keywords'

export const uploadKeywords = async (formData) => {
  const response = await fetch(`${BASE_URL}/${PATH_PREFIX}/upload`, { method: 'POST', body: formData, credentials: "include" })
  const result = await response.json()
  return result
}

export const collectKeywords = async (page, size, sort) => {
  const response = await fetch(`${BASE_URL}/${PATH_PREFIX}?page=${page}&size=${size}&sort=${sort}`, { 
    method: 'GET', credentials: "include" 
  })
  const result = await response.json()
  return result
}