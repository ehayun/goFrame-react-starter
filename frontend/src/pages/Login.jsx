import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

function Login() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [zehut, setZehut] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleLogin = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await fetch('/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ zehut, password }),
        credentials: 'include',
      })

      const data = await response.json()

      if (response.ok) {
        // Redirect to dashboard
        window.location.href = '/dashboard'
      } else {
        setError(data.error || 'Login failed')
      }
    } catch (err) {
      setError('Network error. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  const handleGoogleLogin = () => {
    window.location.href = '/auth/google/login'
  }

  return (
    <div className="min-vh-100 d-flex align-items-center justify-content-center bg-neutral-50">
      <div className="card shadow-lg" style={{ maxWidth: '450px', width: '100%' }}>
        <div className="card-body p-40">
          <div className="text-center mb-32">
            <h3 className="mb-8">ברוך הבא ל-Tzlev</h3>
            <p className="text-secondary-light">התחבר עם תעודת הזהות והסיסמה שלך</p>
          </div>

          {error && (
            <div className="alert alert-danger mb-24" role="alert">
              {error}
            </div>
          )}

          <form onSubmit={handleLogin}>
            <div className="mb-20">
              <label className="form-label fw-semibold text-primary-light text-sm mb-8">
                תעודת זהות
              </label>
              <input
                type="text"
                className="form-control radius-8"
                placeholder="הכנס תעודת זהות"
                value={zehut}
                onChange={(e) => setZehut(e.target.value)}
                required
                maxLength="9"
              />
            </div>

            <div className="mb-20">
              <label className="form-label fw-semibold text-primary-light text-sm mb-8">
                סיסמה
              </label>
              <input
                type="password"
                className="form-control radius-8"
                placeholder="הכנס סיסמה"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>

            <button
              type="submit"
              className="btn btn-primary text-sm btn-sm px-12 py-16 w-100 radius-8 mt-32"
              disabled={loading}
            >
              {loading ? 'מתחבר...' : 'התחבר'}
            </button>
          </form>

          <div className="my-24 text-center">
            <span className="text-secondary-light text-sm">או</span>
          </div>

          <button
            type="button"
            onClick={handleGoogleLogin}
            className="btn btn-outline-primary text-sm btn-sm px-12 py-16 w-100 radius-8 d-flex align-items-center justify-content-center gap-2"
          >
            <iconify-icon icon="logos:google-icon" className="text-xl"></iconify-icon>
            התחבר עם Google
          </button>
        </div>
      </div>
    </div>
  )
}

export default Login
