import React from 'react'
import { useTranslation } from 'react-i18next'

function UserPermissions() {
  const { t } = useTranslation()

  return (
    <div className="tab-pane active">
      <h5>{t('permissions.userPermissions.title')}</h5>
      <p>{t('permissions.userPermissions.description')}</p>
    </div>
  )
}

export default UserPermissions
