import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

function Dashboard() {
  const { t } = useTranslation()
  const [health, setHealth] = useState(null)

  useEffect(() => {
    // Test API call
    fetch('/api/health')
      .then(res => res.json())
      .then(data => setHealth(data))
      .catch(err => console.error('API Error:', err))
  }, [])

  return (
    <div className="d-flex flex-wrap align-items-center justify-content-between gap-3 mb-24">
      <h6 className="fw-semibold mb-0">{t('menu.dashboard')}</h6>
      <div className="card h-100 p-0 radius-12">
        <div className="card-body p-24">
          <h5>Dashboard</h5>
          {health && (
            <div>
              <p>API Status: {health.status}</p>
              <p>Message: {health.message}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Dashboard
