import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import en from '../locales/en/translation.json'
import he from '../locales/he/translation.json'

i18n
  .use(initReactI18next)
  .init({
    resources: {
      en: { translation: en },
      he: { translation: he }
    },
    lng: 'he', // Hebrew as default language
    fallbackLng: 'he',
    interpolation: {
      escapeValue: false
    }
  })

// Set RTL direction for Hebrew by default
document.documentElement.dir = 'rtl'
document.documentElement.lang = 'he'

export default i18n
