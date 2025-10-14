import React from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../context/AuthContext'
import { Link } from 'react-router-dom'

function Navbar() {
  const { t, i18n } = useTranslation()
  const { user, isAuthenticated, logout } = useAuth()

  const changeLanguage = (lng) => {
    i18n.changeLanguage(lng)
    document.documentElement.dir = lng === 'he' ? 'rtl' : 'ltr'
    document.documentElement.lang = lng
  }

  const handleLogout = (e) => {
    e.preventDefault()
    logout()
  }

  return (
    <div className="navbar-header">
      <div className="row align-items-center justify-content-between">
        <div className="col-auto">
          <div className="d-flex flex-wrap align-items-center gap-4">
            <button type="button" className="sidebar-toggle">
              <iconify-icon icon="heroicons:bars-3-solid" className="icon text-2xl non-active"></iconify-icon>
              <iconify-icon icon="iconoir:arrow-right" className="icon text-2xl active"></iconify-icon>
            </button>
            <button type="button" className="sidebar-mobile-toggle">
              <iconify-icon icon="heroicons:bars-3-solid" className="icon"></iconify-icon>
            </button>
          </div>
        </div>
        <div className="col-auto">
          <div className="d-flex flex-wrap align-items-center gap-3">
            {/* Language Switcher - Hidden but kept in code */}
            {/*
            <div className="dropdown">
              <button className="has-indicator w-40-px h-40-px bg-neutral-200 rounded-circle d-flex justify-content-center align-items-center" type="button" data-bs-toggle="dropdown">
                <iconify-icon icon="ph:translate" className="text-primary-light text-xl"></iconify-icon>
              </button>
              <div className="dropdown-menu to-top dropdown-menu-sm">
                <div className="py-12 px-16 radius-8 bg-primary-50 mb-16 d-flex align-items-center justify-content-between gap-2">
                  <div>
                    <h6 className="text-lg text-primary-light fw-semibold mb-0">{t('common.language')}</h6>
                  </div>
                </div>
                <div className="max-h-400-px overflow-y-auto scroll-sm pe-8">
                  <a
                    href="javascript:void(0)"
                    className="dropdown-item d-flex align-items-center gap-2 py-6 px-16"
                    onClick={() => changeLanguage('en')}
                  >
                    <span className="text-secondary-light text-md fw-medium">English</span>
                  </a>
                  <a
                    href="javascript:void(0)"
                    className="dropdown-item d-flex align-items-center gap-2 py-6 px-16"
                    onClick={() => changeLanguage('he')}
                  >
                    <span className="text-secondary-light text-md fw-medium">עברית</span>
                  </a>
                </div>
              </div>
            </div>
            */}

            {/* Theme Toggle */}
            <button type="button" data-theme-toggle className="w-40-px h-40-px bg-neutral-200 rounded-circle d-flex justify-content-center align-items-center">
              <iconify-icon icon="solar:sun-2-bold" className="sun"></iconify-icon>
              <iconify-icon icon="ph:moon-fill" className="moon"></iconify-icon>
            </button>

            {/* User Profile or Login */}
            {isAuthenticated ? (
              <div className="dropdown">
                <button className="d-flex justify-content-center align-items-center rounded-circle" type="button" data-bs-toggle="dropdown">
                  <img src={user?.avatar || "/images/user.png"} alt="User" className="w-40-px h-40-px object-fit-cover rounded-circle" />
                </button>
                <div className="dropdown-menu to-top dropdown-menu-sm">
                  <div className="py-12 px-16 radius-8 bg-primary-50 mb-16 d-flex align-items-center justify-content-between gap-2">
                    <div>
                      <h6 className="text-lg text-primary-light fw-semibold mb-2">
                        {user?.first_name} {user?.last_name}
                      </h6>
                      <span className="text-secondary-light fw-medium text-sm">{user?.email}</span>
                    </div>
                  </div>
                  <ul className="to-top-list">
                    <li>
                      <Link className="dropdown-item text-black px-16 py-8 rounded text-sm" to="/profile">
                        <iconify-icon icon="solar:user-linear" className="icon text-xl"></iconify-icon>
                        {t('common.profile')}
                      </Link>
                    </li>
                    <li>
                      <a className="dropdown-item text-black px-16 py-8 rounded text-sm" href="#" onClick={handleLogout}>
                        <iconify-icon icon="lucide:power" className="icon text-xl"></iconify-icon>
                        {t('common.logout')}
                      </a>
                    </li>
                  </ul>
                </div>
              </div>
            ) : (
              <Link to="/login" className="btn btn-primary text-sm btn-sm px-12 py-10 radius-8">
                התחבר
              </Link>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default Navbar
