import { DashboardLayout } from '@/components/DashboardLayout'
import { ProtectedRoute } from '@/components/ProtectedRoute'

export default function Docs() {
  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div>
          <h1 className="text-3xl font-semibold text-gray-900 mb-2">Documentation</h1>
          <p className="text-gray-600 mb-8">Learn how to integrate Besend with your application</p>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2">
              <div className="card mb-6">
                <h2 className="text-2xl font-semibold text-gray-900 mb-4">Getting Started</h2>
                <div className="space-y-6 text-gray-600">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">1. Create an API Key</h3>
                    <p>Go to the API Keys section and create a new key. You'll use this to authenticate your requests.</p>
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">2. Install SDK (Optional)</h3>
                    <p>We provide SDKs for popular languages, but you can also use our REST API directly.</p>
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">3. Send Your First Email</h3>
                    <p>Make a POST request to our API endpoint with your email details.</p>
                  </div>
                </div>
              </div>

              <div className="card">
                <h2 className="text-2xl font-semibold text-gray-900 mb-4">API Reference</h2>
                <div className="space-y-6">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">Send Email</h3>
                    <p className="text-gray-600 mb-3">POST /api/emails</p>
                    <div className="bg-gray-900 text-gray-100 p-4 rounded-lg text-sm font-mono overflow-x-auto mb-3">
                      <pre>{`curl -X POST https://api.besend.io/api/emails \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "to": "user@example.com",
    "from": "noreply@example.com",
    "subject": "Welcome!",
    "html": "<p>Welcome to Besend</p>"
  }'`}</pre>
                    </div>
                    <div className="bg-blue-50 border border-blue-200 p-4 rounded-lg text-sm text-blue-900">
                      <p className="font-semibold mb-1">Response</p>
                      <code className="text-xs">{`{ "id": "email_123", "status": "sent" }`}</code>
                    </div>
                  </div>

                  <div className="pt-4 border-t border-gray-200">
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">Get Email Status</h3>
                    <p className="text-gray-600 mb-3">GET /api/emails/:id</p>
                    <div className="bg-gray-900 text-gray-100 p-4 rounded-lg text-sm font-mono overflow-x-auto">
                      <pre>{`curl -X GET https://api.besend.io/api/emails/email_123 \\
  -H "Authorization: Bearer YOUR_API_KEY"`}</pre>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <div className="card sticky top-8">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Links</h3>
                <div className="space-y-2">
                  <a href="#" className="flex items-center justify-between p-2.5 rounded-md hover:bg-gray-100 transition-colors">
                    <span className="text-sm font-medium text-gray-900">API Docs</span>
                    <span className="text-gray-500">→</span>
                  </a>
                  <a href="#" className="flex items-center justify-between p-2.5 rounded-md hover:bg-gray-100 transition-colors">
                    <span className="text-sm font-medium text-gray-900">GitHub</span>
                    <span className="text-gray-500">→</span>
                  </a>
                  <a href="#" className="flex items-center justify-between p-2.5 rounded-md hover:bg-gray-100 transition-colors">
                    <span className="text-sm font-medium text-gray-900">Examples</span>
                    <span className="text-gray-500">→</span>
                  </a>
                  <a href="#" className="flex items-center justify-between p-2.5 rounded-md hover:bg-gray-100 transition-colors">
                    <span className="text-sm font-medium text-gray-900">Support</span>
                    <span className="text-gray-500">→</span>
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  )
}
