import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import Modules from '../components/permissions/Modules'
import Groups from '../components/permissions/Groups'
import GroupPermissions from '../components/permissions/GroupPermissions'
import UserPermissions from '../components/permissions/UserPermissions'

function Permissions() {
  const { t } = useTranslation()
  const [activeTab, setActiveTab] = useState('modules')

  const tabs = [
    { id: 'modules', label: t('permissions.tabs.modules') },
    { id: 'groups', label: t('permissions.tabs.groups') },
    { id: 'group-permissions', label: t('permissions.tabs.groupPermissions') },
    { id: 'user-permissions', label: t('permissions.tabs.userPermissions') }
  ]

  const renderTabContent = () => {
    switch (activeTab) {
      case 'modules':
        return <Modules />
      case 'groups':
        return <Groups />
      case 'group-permissions':
        return <GroupPermissions />
      case 'user-permissions':
        return <UserPermissions />
      default:
        return null
    }
  }

  return (
    <div className="container-fluid">
      <div className="row">
        <div className="col-12">
          <div className="page-title-box">
            <h4 className="page-title">{t('permissions.title')}</h4>
          </div>
        </div>
      </div>
      
      <div className="row">
        <div className="col-12">
          <div className="card">
            <div className="card-body">
              {/* Tab Navigation */}
              <ul className="nav nav-tabs" role="tablist">
                {tabs.map((tab) => (
                  <li className="nav-item" key={tab.id}>
                    <button
                      className={`nav-link ${activeTab === tab.id ? 'active' : ''}`}
                      onClick={() => setActiveTab(tab.id)}
                      type="button"
                    >
                      {tab.label}
                    </button>
                  </li>
                ))}
              </ul>

              {/* Tab Content */}
              <div className="tab-content mt-3">
                {renderTabContent()}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Permissions
