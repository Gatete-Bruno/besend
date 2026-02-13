import { DashboardLayout } from '@/components/DashboardLayout'
import { ProtectedRoute } from '@/components/ProtectedRoute'

export default function Usage() {
  return (
    <ProtectedRoute>
      <DashboardLayout>
        <div>
          <h1 className="text-3xl font-semibold text-gray-900 mb-2">Usage & Billing</h1>
          <p className="text-gray-600 mb-8">Track your email usage and manage your billing</p>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2 card">
              <h2 className="text-lg font-semibold text-gray-900 mb-6">Current Plan</h2>
              <div className="space-y-4">
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-sm font-medium text-gray-600">Emails Used This Month</span>
                    <span className="text-sm font-semibold text-gray-900">2,450 / 10,000</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div className="bg-black h-2 rounded-full" style={{ width: '24.5%' }}></div>
                  </div>
                </div>

                <div className="pt-4 border-t border-gray-200">
                  <p className="text-sm text-gray-600 mb-4">
                    You're on the <strong>Professional</strong> plan. Upgrade or downgrade anytime.
                  </p>
                  <button className="btn-primary">
                    Manage Billing
                  </button>
                </div>
              </div>
            </div>

            <div className="card">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">Billing Info</h2>
              <div className="space-y-3 text-sm">
                <div>
                  <p className="text-gray-600">Plan</p>
                  <p className="font-semibold text-gray-900">Professional</p>
                </div>
                <div>
                  <p className="text-gray-600">Billing Cycle</p>
                  <p className="font-semibold text-gray-900">Monthly</p>
                </div>
                <div>
                  <p className="text-gray-600">Next Invoice</p>
                  <p className="font-semibold text-gray-900">March 1, 2024</p>
                </div>
              </div>
            </div>
          </div>

          <div className="card mt-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Usage History</h2>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-200">
                    <th className="text-left py-3 px-0 font-semibold text-gray-900">Date</th>
                    <th className="text-right py-3 px-0 font-semibold text-gray-900">Emails Sent</th>
                    <th className="text-right py-3 px-0 font-semibold text-gray-900">Cost</th>
                  </tr>
                </thead>
                <tbody>
                  <tr className="border-b border-gray-200">
                    <td className="py-3 px-0 text-gray-600">February 2024</td>
                    <td className="text-right py-3 px-0 text-gray-900">8,250</td>
                    <td className="text-right py-3 px-0 text-gray-900">$29.00</td>
                  </tr>
                  <tr className="border-b border-gray-200">
                    <td className="py-3 px-0 text-gray-600">January 2024</td>
                    <td className="text-right py-3 px-0 text-gray-900">5,120</td>
                    <td className="text-right py-3 px-0 text-gray-900">$29.00</td>
                  </tr>
                  <tr>
                    <td className="py-3 px-0 text-gray-600">December 2023</td>
                    <td className="text-right py-3 px-0 text-gray-900">3,890</td>
                    <td className="text-right py-3 px-0 text-gray-900">$0.00</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  )
}
