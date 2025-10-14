import React from 'react'
import { useTranslation } from 'react-i18next'

function Home() {
  const { t } = useTranslation()

  return (
    <div className="d-flex flex-wrap align-items-center justify-content-between gap-3 mb-24">
      <h6 className="fw-semibold mb-0">{t('menu.home')}</h6>
      <div className="card h-100 p-0 radius-12">
        <div className="card-body p-24">
          <h5>Welcome to Tzlev!</h5>
          <p>This is the home page indeed.</p>
        </div>
      </div>
    </div>
  )
}

export default Home
