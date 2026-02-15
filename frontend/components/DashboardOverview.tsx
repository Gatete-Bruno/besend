// ============================================
// IMPORTS
// ============================================
import React, { useState, useEffect } from 'react'
import axios from 'axios'
import { useAuth } from '@/context/AuthContext'
import Link from 'next/link'

// ============================================
// TYPES
// ============================================
interface DashboardStats {
  emailsSentToday: number
  emailsSentMonth: number
  totalEmails: number
  apiKeys: number
}

// ============================================
// COMPONENT
// ============================================
export const DashboardOverview: React.FC = () => {
  // ============================================
  // STATE
  // ============================================
  const { user, token } = useAuth()
  const [stats, setStats] = useState<DashboardStats>({
    emailsSentToday: 0,
    emailsSentMonth: 0,
    totalEmails: 0,
    apiKeys: 0,
  })
  const [loading, setLoading] = useState(true)

  // ============================================
  // CONSTANTS
  // ============================================
  const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

  // ============================================
  // EFFECTS
  // ============================================
  useEffect(() => {
    fetchStats()
  }, [])

  // ============================================
  // FUNCTIONS
  // ============================================
  const fetchStats = async () => {
    try {
      const response = await axios.get(`${API_URL}/api/stats`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      setStats(response.data)
    } catch (error) {
      console.error('Failed to fetch stats:', error)
    } finally {
      setLoading(false)
    }
  }

  // ============================================
  // RENDER
  // ============================================
  return (
    <div>
      {/* ============================================ */}
      {/* HEADER SECTION */}
      {/* ============================================ */}
      <div className="mb-8">
        <h1 className="text-3xl font-semibold text-white mb-2">
          Welcome back, {user?.name?.split(' ')[0]}
        </h1>
        <p className="text-gray-400">Track your email sending activity and manage your account</p>
      </div>

      {/* ============================================ */}
      {/* LOADING STATE */}
      {/* ============================================ */}
      {loading ? (
        <div className="text-center py-12">
          <div className="w-6 h-6 border-2 border-gray-700 border-t-white rounded-full animate-spin mx-auto"></div>
        </div>
      ) : (
        <>
          {/* ============================================ */}
          {/* STATS CARDS SECTION */}
          {/* ============================================ */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            {/* Card 1: Sent Today */}
            <div className="card">
              <p className="text-gray-400 text-sm mb-1">Sent Today</p>
              <p className="text-3xl font-semibold text-white">
                {stats.emailsSentToday.toLocaleString()}
              </p>
            </div>

            {/* Card 2: Sent This Month */}
            <div className="card">
              <p className="text-gray-400 text-sm mb-1">Sent This Month</p>
              <p className="text-3xl font-semibold text-white">
                {stats.emailsSentMonth.toLocaleString()}
              </p>
            </div>

            {/* Card 3: Total Emails */}
            <div className="card">
              <p className="text-gray-400 text-sm mb-1">Total Emails</p>
              <p className="text-3xl font-semibold text-white">
                {stats.totalEmails.toLocaleString()}
              </p>
            </div>

            {/* Card 4: API Keys */}
            <div className="card">
              <p className="text-gray-400 text-sm mb-1">API Keys</p>
              <p className="text-3xl font-semibold text-white">{stats.apiKeys}</p>
            </div>
          </div>

          {/* ============================================ */}
          {/* MAIN CONTENT GRID */}
          {/* ============================================ */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* ============================================ */}
            {/* LEFT COLUMN: GETTING STARTED GUIDE */}
            {/* ============================================ */}
            <div className="lg:col-span-2 card">
              <h2 className="text-lg font-semibold text-white mb-6">Getting Started</h2>
              <div className="space-y-6">
                {/* Step 1: Create API Key */}
                <div className="flex gap-4">
                  <div className="w-6 h-6 rounded-full bg-white text-black flex items-center justify-center text-sm font-semibold flex-shrink-0">
                    1
                  </div>
                  <div>
                    <h3 className="font-semibold text-white mb-1">Create an API Key</h3>
                    <p className="text-sm text-gray-400 mb-3">
                      Generate your first API key to start integrating with Besend.
                    </p>
                    <Link href="/dashboard/api-keys" className="text-sm font-medium text-white hover:text-gray-300">
                      Go to API Keys →
                    </Link>
                  </div>
                </div>

                {/* Divider */}
                <div className="divider"></div>

                {/* Step 2: Read Documentation */}
                <div className="flex gap-4">
                  <div className="w-6 h-6 rounded-full bg-white text-black flex items-center justify-center text-sm font-semibold flex-shrink-0">
                    2
                  </div>
                  <div>
                    <h3 className="font-semibold text-white mb-1">Read Documentation</h3>
                    <p className="text-sm text-gray-400 mb-3">
                      Learn how to send transactional emails using our REST API.
                    </p>
                    <Link href="/dashboard/docs" className="text-sm font-medium text-white hover:text-gray-300">
                      View Documentation →
                    </Link>
                  </div>
                </div>

                {/* Divider */}
                <div className="divider"></div>

                {/* Step 3: Send First Email */}
                <div className="flex gap-4">
                  <div className="w-6 h-6 rounded-full bg-white text-black flex items-center justify-center text-sm font-semibold flex-shrink-0">
                    3
                  </div>
                  <div>
                    <h3 className="font-semibold text-white mb-1">Send First Email</h3>
                    <p className="text-sm text-gray-400 mb-3">
                      Test your integration with a simple API request.
                    </p>
                    <a href="#" className="text-sm font-medium text-white hover:text-gray-300">
                      Send Test Email →
                    </a>
                  </div>
                </div>
              </div>
            </div>

            {/* ============================================ */}
            {/* RIGHT COLUMN: RESOURCES */}
            {/* ============================================ */}
            <div className="card">
              <h2 className="text-lg font-semibold text-white mb-4">Resources</h2>
              <div className="space-y-2">
                {/* Resource 1: API Docs */}
                <a href="https://docs.besend.io" target="_blank" rel="noopener noreferrer" className="flex items-center justify-between p-2.5 rounded-md hover:bg-gray-800 transition-colors">
                  <span className="text-sm font-medium text-white">API Docs</span>
                  <span className="text-gray-500">→</span>
                </a>

                {/* Resource 2: Status Page */}
                <a href="https://status.besend.io" target="_blank" rel="noopener noreferrer" className="flex items-center justify-between p-2.5 rounded-md hover:bg-gray-800 transition-colors">
                  <span className="text-sm font-medium text-white">Status Page</span>
                  <span className="text-gray-500">→</span>
                </a>

                {/* Resource 3: Support Email */}
                <a href="mailto:support@besend.io" className="flex items-center justify-between p-2.5 rounded-md hover:bg-gray-800 transition-colors">
                  <span className="text-sm font-medium text-white">Support</span>
                  <span className="text-gray-500">→</span>
                </a>
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  )
}
