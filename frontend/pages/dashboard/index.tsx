import { DashboardLayout } from '@/components/DashboardLayout'
import { DashboardOverview } from '@/components/DashboardOverview'
import { ProtectedRoute } from '@/components/ProtectedRoute'

export default function Dashboard() {
  return (
    <ProtectedRoute>
      <DashboardLayout>
        <DashboardOverview />
      </DashboardLayout>
    </ProtectedRoute>
  )
}
