import React, { useState, useEffect } from 'react'
import axios from 'axios'
import { useAuth } from '@/context/AuthContext'

interface APIKey {
  id: string
  name: string
  key: string
  createdAt: string
  lastUsed: string | null
  isActive: boolean
}

export const APIKeysPage: React.FC = () => {
  const { token } = useAuth()
  const [keys, setKeys] = useState<APIKey[]>([])
  const [loading, setLoading] = useState(true)
  const [creating, setCreating] = useState(false)
  const [keyName, setKeyName] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [error, setError] = useState('')
  const [copiedId, setCopiedId] = useState<string | null>(null)

  const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

  useEffect(() => {
    fetchKeys()
  }, [])

  const fetchKeys = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/keys`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      setKeys(response.data || [])
    } catch (error) {
      console.error('Failed to fetch keys:', error)
      setError('Failed to load API keys')
    } finally {
      setLoading(false)
    }
  }

  const handleCreateKey = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setCreating(true)

    try {
      const response = await axios.post(
        `${API_URL}/api/keys`,
        { name: keyName },
        { headers: { Authorization: `Bearer ${token}` } }
      )
      setKeys([response.data, ...keys])
      setKeyName('')
      setShowForm(false)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create key')
    } finally {
      setCreating(false)
    }
  }

  const handleRevokeKey = async (keyId: string) => {
    if (!confirm('Are you sure you want to revoke this key?')) return

    try {
      await axios.delete(`${API_URL}/api/keys/${keyId}`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      setKeys(keys.filter((k) => k.id !== keyId))
    } catch (error) {
      setError('Failed to revoke key')
    }
  }

  const copyToClipboard = (key: string, id: string) => {
    navigator.clipboard.writeText(key)
    setCopiedId(id)
    setTimeout(() => setCopiedId(null), 2000)
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-semibold text-gray-900 mb-2">API Keys</h1>
          <p className="text-gray-600">Manage your API keys for Besend integration</p>
        </div>
        {!showForm && (
          <button
            onClick={() => setShowForm(true)}
            className="btn-primary"
          >
            + Create Key
          </button>
        )}
      </div>

      {error && (
        <div className="p-4 bg-red-50 border border-red-200 rounded-lg mb-6 text-red-800 text-sm">
          {error}
        </div>
      )}

      {showForm && (
        <div className="card mb-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Create New API Key</h2>
          <form onSubmit={handleCreateKey} className="space-y-4">
            <div>
              <label htmlFor="keyName" className="block text-sm font-medium text-gray-900 mb-1.5">
                Key Name
              </label>
              <input
                id="keyName"
                type="text"
                value={keyName}
                onChange={(e) => setKeyName(e.target.value)}
                className="input-field"
                placeholder="Production API Key"
                required
              />
            </div>
            <div className="flex gap-2">
              <button type="submit" disabled={creating} className="btn-primary">
                {creating ? 'Creating...' : 'Create Key'}
              </button>
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="btn-secondary"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <div className="text-center py-12">
          <div className="w-6 h-6 border-2 border-gray-300 border-t-black rounded-full animate-spin mx-auto"></div>
        </div>
      ) : keys.length === 0 ? (
        <div className="card text-center py-12">
          <p className="text-gray-600 mb-4">No API keys yet</p>
          <button
            onClick={() => setShowForm(true)}
            className="btn-primary"
          >
            Create Your First Key
          </button>
        </div>
      ) : (
        <div className="space-y-3">
          {keys.map((key) => (
            <div key={key.id} className="card">
              <div className="flex items-center justify-between">
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900 mb-2">{key.name}</h3>
                  <div className="flex items-center gap-2 mb-2">
                    <code className="text-sm bg-gray-100 px-2.5 py-1.5 rounded font-mono text-gray-900 flex-1 overflow-hidden text-ellipsis">
                      {key.key}
                    </code>
                    <button
                      onClick={() => copyToClipboard(key.key, key.id)}
                      className="px-2.5 py-1.5 text-sm bg-gray-100 hover:bg-gray-200 rounded transition-colors"
                    >
                      {copiedId === key.id ? '✓ Copied' : 'Copy'}
                    </button>
                  </div>
                  <p className="text-xs text-gray-500">
                    Created {new Date(key.createdAt).toLocaleDateString()}
                    {key.lastUsed && ` • Last used ${new Date(key.lastUsed).toLocaleDateString()}`}
                  </p>
                </div>
                <button
                  onClick={() => handleRevokeKey(key.id)}
                  className="ml-4 px-3 py-1.5 text-sm text-red-600 hover:bg-red-50 rounded transition-colors border border-red-200"
                >
                  Revoke
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="card mt-8 bg-gray-50">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Example Usage</h2>
        <div className="bg-gray-900 text-gray-100 p-4 rounded-lg text-sm font-mono overflow-x-auto">
          <pre>{`curl -X POST https://api.besend.io/api/emails \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "to": "user@example.com",
    "from": "noreply@example.com",
    "subject": "Hello",
    "html": "<p>Welcome!</p>"
  }'`}</pre>
        </div>
      </div>
    </div>
  )
}
