import { DashboardLayout } from '@/components/DashboardLayout'
import { APIKeysPage } from '@/components/APIKeysPage'
import { ProtectedRoute } from '@/components/ProtectedRoute'

export default function ApiKeys() {
  return (
    <ProtectedRoute>
      <DashboardLayout>
        <APIKeysPage />
      </DashboardLayout>
    </ProtectedRoute>
  )
}
