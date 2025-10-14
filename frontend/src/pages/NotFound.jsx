import React from 'react'
import { Link } from 'react-router-dom'

function NotFound() {
  return (
    <div className="d-flex flex-wrap align-items-center justify-content-between gap-3 mb-24">
      <div className="card h-100 p-0 radius-12 w-100">
        <div className="card-body p-40 text-center">
          <div className="mb-32">
            <iconify-icon icon="mdi:alert-circle-outline" className="text-danger-main" style={{ fontSize: '120px' }}></iconify-icon>
          </div>
          <h1 className="mb-16" style={{ fontSize: '72px', fontWeight: '700' }}>404</h1>
          <h4 className="mb-16">הדף לא נמצא</h4>
          <p className="text-secondary-light mb-32">
            מצטערים, הדף שחיפשת לא קיים במערכת.
          </p>
          <Link to="/" className="btn btn-primary text-sm btn-sm px-24 py-12 radius-8">
            <iconify-icon icon="solar:home-2-linear" className="icon text-lg me-8"></iconify-icon>
            חזרה לדף הבית
          </Link>
        </div>
      </div>
    </div>
  )
}

export default NotFound
