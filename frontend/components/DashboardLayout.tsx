import React from 'react'
import Link from 'next/link'
import { useRouter } from 'next/router'
import { useAuth } from '@/context/AuthContext'

interface DashboardLayoutProps {
  children: React.ReactNode
}

export const DashboardLayout: React.FC<DashboardLayoutProps> = ({ children }) => {
  const { user, logout } = useAuth()
  const router = useRouter()

  const menuItems = [
    { href: '/dashboard', label: 'Overview', icon: '⊞' },
    { href: '/dashboard/api-keys', label: 'API Keys', icon: '🔑' },
    { href: '/dashboard/usage', label: 'Usage & Billing', icon: '📊' },
    { href: '/dashboard/docs', label: 'Documentation', icon: '📖' },
  ]

  const isActive = (href: string) => router.pathname === href

  return (
    <div className="flex h-screen bg-white">
      <aside className="w-64 bg-white border-r border-gray-200 flex flex-col">
        <div className="p-6 border-b border-gray-200">
          <Link href="/dashboard" className="flex items-center gap-2">
            <div className="w-8 h-8 bg-black rounded-md flex items-center justify-center font-semibold text-white text-sm">
              B
            </div>
            <span className="font-semibold text-lg">besend</span>
          </Link>
        </div>

        <nav className="flex-1 p-4 space-y-1">
          {menuItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-md transition-colors text-sm font-medium ${
                isActive(item.href)
                  ? 'bg-gray-100 text-black'
                  : 'text-gray-700 hover:bg-gray-50'
              }`}
            >
              <span>{item.icon}</span>
              {item.label}
            </Link>
          ))}
        </nav>

        <div className="p-4 border-t border-gray-200 space-y-3">
          <div className="px-3 py-2.5 bg-gray-50 rounded-md">
            <p className="text-xs text-gray-600 mb-0.5">Signed in as</p>
            <p className="text-sm font-medium truncate text-gray-900">{user?.email}</p>
          </div>
          <button
            onClick={() => {
              logout()
              router.push('/login')
            }}
            className="w-full btn-secondary text-sm"
          >
            Sign Out
          </button>
        </div>
      </aside>

      <main className="flex-1 overflow-auto bg-gray-50">
        <div className="p-8">
          {children}
        </div>
      </main>
    </div>
  )
}
