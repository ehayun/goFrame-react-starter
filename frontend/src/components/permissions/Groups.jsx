import React from 'react'
import { useTranslation } from 'react-i18next'

function Groups() {
  const { t } = useTranslation()

  return (
    <div className="tab-pane active">
      <h5>{t('permissions.groups.title')}</h5>
      <p>{t('permissions.groups.description')}</p>
    </div>
  )
}

export default Groups
