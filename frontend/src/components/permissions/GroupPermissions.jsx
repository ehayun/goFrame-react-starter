import React from 'react'
import { useTranslation } from 'react-i18next'

function GroupPermissions() {
  const { t } = useTranslation()

  return (
    <div className="tab-pane active">
      <h5>{t('permissions.groupPermissions.title')}</h5>
      <p>{t('permissions.groupPermissions.description')}</p>
    </div>
  )
}

export default GroupPermissions
