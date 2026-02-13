import React, { createContext, useContext, useState, useEffect } from 'react'
import axios from 'axios'
import Cookies from 'js-cookie'

interface User {
  id: string
  email: string
  name: string
  createdAt: string
}

interface AuthContextType {
  user: User | null
  token: string | null
  loading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string) => Promise<void>
  logout: () => void
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

  useEffect(() => {
    const storedToken = Cookies.get('auth_token')
    if (storedToken) {
      setToken(storedToken)
      fetchUser(storedToken)
    } else {
      setLoading(false)
    }
  }, [])

  const fetchUser = async (authToken: string) => {
    try {
      const response = await axios.get(`${API_URL}/api/users/me`, {
        headers: { Authorization: `Bearer ${authToken}` },
      })
      setUser(response.data)
    } catch (error) {
      console.error('Failed to fetch user:', error)
      Cookies.remove('auth_token')
      setToken(null)
    } finally {
      setLoading(false)
    }
  }

  const login = async (email: string, password: string) => {
    const response = await axios.post(`${API_URL}/api/auth/login`, { email, password })
    const { token: newToken, user: newUser } = response.data
    Cookies.set('auth_token', newToken, { expires: 7 })
    setToken(newToken)
    setUser(newUser)
  }

  const register = async (email: string, password: string, name: string) => {
    const response = await axios.post(`${API_URL}/api/auth/register`, { email, password, name })
    const { token: newToken, user: newUser } = response.data
    Cookies.set('auth_token', newToken, { expires: 7 })
    setToken(newToken)
    setUser(newUser)
  }

  const logout = () => {
    Cookies.remove('auth_token')
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, token, loading, login, register, logout, isAuthenticated: !!user }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
