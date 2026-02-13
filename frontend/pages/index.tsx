import React from 'react'
import Link from 'next/link'
import { useRouter } from 'next/router'
import { useAuth } from '@/context/AuthContext'

export default function Home() {
  const { isAuthenticated, loading } = useAuth()
  const router = useRouter()

  if (loading) {
    return null
  }

  if (isAuthenticated) {
    if (!router.pathname.includes('/dashboard')) {
      router.push('/dashboard')
    }
    return null
  }

  return (
    <div className="min-h-screen bg-black">
      <nav className="border-b border-gray-800 bg-black sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2 group">
            <div className="w-8 h-8 bg-white rounded-md flex items-center justify-center font-semibold text-black text-sm">
              B
            </div>
            <span className="font-semibold text-lg text-white">besend</span>
          </Link>
          <div className="flex items-center gap-4">
            <Link href="/login" className="text-gray-300 hover:text-white transition-colors font-medium">
              Sign In
            </Link>
            <Link href="/register" className="btn-primary">
              Get Started
            </Link>
          </div>
        </div>
      </nav>

      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 sm:py-32">
        <div className="text-center mb-16">
          <h1 className="text-5xl sm:text-6xl font-semibold text-white mb-6 leading-tight">
            Email infrastructure for developers
          </h1>
          <p className="text-xl text-gray-400 mb-8 max-w-2xl mx-auto">
            Send transactional emails at scale with a simple REST API. Built for reliability, deliverability, and developer experience.
          </p>
          <div className="flex items-center justify-center gap-4 flex-wrap">
            <Link href="/register" className="btn-primary">
              Start for free
            </Link>
            <a href="#features" className="btn-secondary">
              Learn more
            </a>
          </div>
        </div>
      </section>

      <section id="features" className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <h2 className="text-4xl font-semibold text-white mb-12 text-center">Features</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="card">
            <div className="text-4xl mb-4">🚀</div>
            <h3 className="text-lg font-semibold text-white mb-2">Simple API</h3>
            <p className="text-gray-400">Easy-to-use REST API for sending emails in seconds</p>
          </div>
          <div className="card">
            <div className="text-4xl mb-4">📊</div>
            <h3 className="text-lg font-semibold text-white mb-2">Real-time Analytics</h3>
            <p className="text-gray-400">Track deliverability and engagement with detailed metrics</p>
          </div>
          <div className="card">
            <div className="text-4xl mb-4">🔒</div>
            <h3 className="text-lg font-semibold text-white mb-2">Secure & Reliable</h3>
            <p className="text-gray-400">Enterprise-grade infrastructure with 99.9% uptime SLA</p>
          </div>
        </div>
      </section>

      <section className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <h2 className="text-4xl font-semibold text-white mb-12 text-center">Pricing</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="card">
            <h3 className="text-lg font-semibold text-white mb-2">Starter</h3>
            <p className="text-3xl font-bold text-white mb-1">Free</p>
            <p className="text-gray-400 text-sm mb-6">Forever free for up to 100 emails/day</p>
            <button className="w-full btn-secondary mb-6">Get started</button>
            <ul className="space-y-2 text-sm text-gray-400">
              <li>✓ 100 emails/day</li>
              <li>✓ Basic analytics</li>
              <li>✓ Email support</li>
            </ul>
          </div>
          <div className="card border-2 border-white">
            <h3 className="text-lg font-semibold text-white mb-2">Professional</h3>
            <p className="text-3xl font-bold text-white mb-1">$29</p>
            <p className="text-gray-400 text-sm mb-6">Per month, billed annually</p>
            <button className="w-full btn-primary mb-6">Start free trial</button>
            <ul className="space-y-2 text-sm text-gray-400">
              <li>✓ 10,000 emails/month</li>
              <li>✓ Advanced analytics</li>
              <li>✓ Priority support</li>
            </ul>
          </div>
          <div className="card">
            <h3 className="text-lg font-semibold text-white mb-2">Enterprise</h3>
            <p className="text-3xl font-bold text-white mb-1">Custom</p>
            <p className="text-gray-400 text-sm mb-6">For high-volume senders</p>
            <button className="w-full btn-secondary mb-6">Contact us</button>
            <ul className="space-y-2 text-sm text-gray-400">
              <li>✓ Unlimited emails</li>
              <li>✓ Custom integration</li>
              <li>✓ Dedicated support</li>
            </ul>
          </div>
        </div>
      </section>

      <footer className="border-t border-gray-800 bg-black">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
            <div>
              <h4 className="font-semibold text-white mb-4">Product</h4>
              <ul className="space-y-2 text-sm text-gray-400">
                <li><a href="#" className="hover:text-white">Features</a></li>
                <li><a href="#" className="hover:text-white">Pricing</a></li>
                <li><a href="#" className="hover:text-white">API Docs</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-white mb-4">Company</h4>
              <ul className="space-y-2 text-sm text-gray-400">
                <li><a href="#" className="hover:text-white">About</a></li>
                <li><a href="#" className="hover:text-white">Blog</a></li>
                <li><a href="#" className="hover:text-white">Contact</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-white mb-4">Legal</h4>
              <ul className="space-y-2 text-sm text-gray-400">
                <li><a href="#" className="hover:text-white">Privacy</a></li>
                <li><a href="#" className="hover:text-white">Terms</a></li>
                <li><a href="#" className="hover:text-white">Security</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-white mb-4">Social</h4>
              <ul className="space-y-2 text-sm text-gray-400">
                <li><a href="#" className="hover:text-white">Twitter</a></li>
                <li><a href="#" className="hover:text-white">GitHub</a></li>
                <li><a href="#" className="hover:text-white">LinkedIn</a></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 pt-8 text-center text-sm text-gray-400">
            <p>&copy; 2024 Besend. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
