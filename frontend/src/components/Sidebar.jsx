import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { menuConfig } from '../config/menuConfig'

function Sidebar() {
  const { t } = useTranslation()
  const location = useLocation()

  const isActive = (path) => location.pathname === path

  return (
    <aside className="sidebar">
      <button type="button" className="sidebar-close-btn">
        <iconify-icon icon="radix-icons:cross-2"></iconify-icon>
      </button>
      <div>
        <a href="/" className="sidebar-logo">
          <img src="/images/logo.png" alt="site logo" className="light-logo" />
          <img src="/images/logo-light.png" alt="site logo" className="dark-logo" />
          <img src="/images/logo-icon.png" alt="site logo" className="logo-icon" />
        </a>
      </div>
      <div className="sidebar-menu-area">
        <ul className="sidebar-menu" id="sidebar-menu">
          {menuConfig.map((item) => (
            <li key={item.id}>
              <Link
                to={item.path}
                className={isActive(item.path) ? 'active-page' : ''}
              >
                <iconify-icon icon={item.icon} className="menu-icon"></iconify-icon>
                <span className={"m e-2"}>{t(item.titleKey)}</span>
              </Link>
            </li>
          ))}
        </ul>
      </div>
    </aside>
  )
}

export default Sidebar
