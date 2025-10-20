import React, { useEffect } from 'react'
import { Routes, Route } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import Sidebar from './components/Sidebar'
import Navbar from './components/Navbar'
import Footer from './components/Footer'
import ThemeCustomization from './components/ThemeCustomization'
import Home from './pages/Home'
import Dashboard from './pages/Dashboard'
import Login from './pages/Login'
import Permissions from './pages/Permissions'
import NotFound from './pages/NotFound'

// Layout wrapper for all pages
function MainLayout({ children }) {
  // Load template JavaScript after React mounts
  useEffect(() => {
    const script = document.createElement('script')
    script.src = '/js/app.js'
    script.async = false
    document.body.appendChild(script)

    return () => {
      if (document.body.contains(script)) {
        document.body.removeChild(script)
      }
    }
  }, [])

  return (
    <>
      <style>{`
        .academic-years-wrapper {
          padding: 0 15px;
        }
        
        .academic-years-container {
          width: 100%;
        }
        
        .academic-years-select {
          width: 100%;
        }
        
        .academic-years-loading {
          display: flex;
          justify-content: center;
          align-items: center;
          padding: 8px;
        }
        
        .spinner-border-sm {
          width: 1rem;
          height: 1rem;
          border-width: 0.125em;
        }
        
        .spinner-border {
          display: inline-block;
          width: 2rem;
          height: 2rem;
          vertical-align: -0.125em;
          border: 0.25em solid currentColor;
          border-right-color: transparent;
          border-radius: 50%;
          animation: spinner-border 0.75s linear infinite;
        }
        
        @keyframes spinner-border {
          to {
            transform: rotate(360deg);
          }
        }
        
        .visually-hidden {
          position: absolute !important;
          width: 1px !important;
          height: 1px !important;
          padding: 0 !important;
          margin: -1px !important;
          overflow: hidden !important;
          clip: rect(0, 0, 0, 0) !important;
          white-space: nowrap !important;
          border: 0 !important;
        }
        
        .navbar-academic-year {
          min-width: 200px;
        }
        
        .navbar-academic-year .academic-years-container {
          width: 100%;
        }
        
        .navbar-academic-year .academic-years-select {
          width: 100%;
        }
        
        .navbar-academic-year .academic-years-select .react-select__control {
          min-height: 36px;
          font-size: 14px;
        }
      `}</style>
      <ThemeCustomization />
      <Sidebar />
      <main className="dashboard-main">
        <Navbar />
        <div className="dashboard-main-body">
          {children}
        </div>
        <Footer />
      </main>
    </>
  )
}

function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route
        path="/"
        element={
          <MainLayout>
            <Home />
          </MainLayout>
        }
      />
      <Route
        path="/dashboard"
        element={
          <MainLayout>
            <Dashboard />
          </MainLayout>
        }
      />
      <Route
        path="/permissions"
        element={
          <MainLayout>
            <Permissions />
          </MainLayout>
        }
      />
      {/* Catch-all route for 404 - must be last */}
      <Route
        path="*"
        element={
          <MainLayout>
            <NotFound />
          </MainLayout>
        }
      />
    </Routes>
  )
}

function App() {
  return (
    <AuthProvider>
      <AppRoutes />
    </AuthProvider>
  )
}

export default App
