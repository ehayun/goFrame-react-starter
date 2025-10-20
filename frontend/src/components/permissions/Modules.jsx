import React, {useEffect, useState} from 'react'
import {useTranslation} from 'react-i18next'
import {useAuth} from '../../context/AuthContext'

function Modules() {
    const {t} = useTranslation()
    const {isAuthenticated, loading: authLoading} = useAuth()
    const [resources, setResources] = useState([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)
    const [showModal, setShowModal] = useState(false)
    const [editingResource, setEditingResource] = useState(null)
    const [formData, setFormData] = useState({name: '', description: ''})

    useEffect(() => {
        // Wait for auth check to complete
        if (authLoading) {
            return
        }

        if (isAuthenticated) {
            loadResources()
        } else {
            setLoading(false)
            setError(t('permissions.modules.authRequired'))
        }
    }, [isAuthenticated, authLoading])

    const loadResources = async () => {
        try {
            setLoading(true)
            setError(null)

            const response = await fetch('/api/app-resources', {
                method: 'GET',
                credentials: 'include'
            })

            if (response.ok) {
                const data = await response.json()
                setResources(data.resources || [])
            } else {
                const errorData = await response.json()
                console.error('Load resources error response:', errorData)
                console.error('Response status:', response.status)
                setError(errorData.error || errorData.message || t('permissions.modules.loadError'))
            }
        } catch (err) {
            setError(t('permissions.modules.loadError'))
            console.error('Error loading resources:', err)
        } finally {
            setLoading(false)
        }
    }

    const handleCreate = () => {
        setEditingResource(null)
        setFormData({name: '', description: ''})
        setShowModal(true)
    }

    const handleEdit = (resource) => {
        setEditingResource(resource)
        setFormData({name: resource.name, description: resource.description || ''})
        setShowModal(true)
    }

    const handleDelete = async (id) => {
        if (!window.confirm(t('permissions.modules.confirmDelete'))) {
            return
        }

        // Check authentication before attempting delete
        if (!isAuthenticated) {
            setError(t('permissions.modules.authRequired'))
            return
        }

        // Validate id
        if (!id || id === '') {
            setError('Invalid resource ID')
            return
        }

        const url = `/api/app-resources/${encodeURIComponent(id)}`
        try {
            const response = await fetch(url, {
                method: 'DELETE',
                credentials: 'include'
            })

            if (response.ok) {
                await loadResources()
            } else {
                const errorData = await response.json()
                console.error('Delete error response:', errorData)
                console.error('Response status:', response.status)
                console.error('Response headers:', Object.fromEntries(response.headers.entries()))

                // Handle specific error cases
                if (response.status === 401) {
                    // Try to refresh authentication
                    try {
                        const authResponse = await fetch('/api/me', {
                            credentials: 'include'
                        })
                        if (authResponse.ok) {
                            // Session is valid, retry the delete
                            const retryResponse = await fetch(`/api/app-resources/${id}`, {
                                method: 'DELETE',
                                credentials: 'include'
                            })
                            if (retryResponse.ok) {
                                await loadResources()
                                return
                            }
                        }
                    } catch (authErr) {
                        console.error('Auth refresh failed:', authErr)
                    }
                    setError(t('permissions.modules.authRequired'))
                } else {
                    setError(errorData.error || errorData.message || t('permissions.modules.deleteError'))
                }
            }
        } catch (err) {
            setError(t('permissions.modules.deleteError'))
            console.error('Error deleting resource:', err)
        }
    }

    const handleSubmit = async (e) => {
        e.preventDefault()

        if (!formData.name.trim()) {
            setError(t('permissions.modules.nameRequired'))
            return
        }

        try {
            const url = editingResource
                ? `/api/app-resources/${editingResource.id}`
                : '/api/app-resources'

            const method = editingResource ? 'PUT' : 'POST'

            const response = await fetch(url, {
                method,
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify(formData)
            })

            if (response.ok) {
                setShowModal(false)
                await loadResources()
            } else {
                const errorData = await response.json()
                setError(errorData.error || t('permissions.modules.saveError'))
            }
        } catch (err) {
            setError(t('permissions.modules.saveError'))
            console.error('Error saving resource:', err)
        }
    }

    const handleCloseModal = () => {
        setShowModal(false)
        setEditingResource(null)
        setFormData({name: '', description: ''})
        setError(null)
    }

    if (loading) {
        return (
            <div className="tab-pane active">
                <h5>{t('permissions.modules.title')}</h5>
                <div className="d-flex justify-content-center">
                    <div className="spinner-border" role="status">
                        <span className="visually-hidden">Loading...</span>
                    </div>
                </div>
            </div>
        )
    }

    if (error) {
        return (
            <div className="tab-pane active">
                <h5>{t('permissions.modules.title')}</h5>
                <div className="alert alert-danger" role="alert">
                    {error}
                </div>
            </div>
        )
    }

    return (
        <div className="tab-pane active">
            <div className="d-flex justify-content-between align-items-center mb-3">
                <h5>{t('permissions.modules.title')}</h5>
                <button
                    className="btn btn-primary btn-sm"
                    onClick={handleCreate}
                >
                    <i className="fas fa-plus me-1"></i>
                    {t('permissions.modules.addNew')}
                </button>
            </div>

            {error && (
                <div className="alert alert-danger" role="alert">
                    {error}
                </div>
            )}

            <div className="table-responsive">
                <table className="table table-striped table-sm">
                    <thead>
                    <tr>
                        <th className="text-start">{t('permissions.modules.name')}</th>
                        <th className="text-start">{t('permissions.modules.description')}</th>
                        <th className="text-center">{t('permissions.modules.actions')}</th>
                    </tr>
                    </thead>
                    <tbody>
                    {resources.length === 0 ? (
                        <tr>
                            <td colSpan="3" className="text-center text-muted">
                                {t('permissions.modules.noResources')}
                            </td>
                        </tr>
                    ) : (
                        resources.map((resource) => {
                            return (
                                <tr key={resource.id}>
                                    <td>{resource.name}</td>
                                    <td>{resource.description || '-'}</td>
                                    <td className="text-center">
                                        <button
                                            className="btn btn-sm btn-outline-primary me-1"
                                            onClick={() => handleEdit(resource)}
                                            title={t('permissions.modules.edit')}
                                        >
                                            <iconify-icon icon="solar:pen-bold" style={{verticalAlign: 'middle'}}>
                                                <span className="ms-1">{t('permissions.modules.edit')}</span>
                                            </iconify-icon>

                                        </button>
                                        <button
                                            className="btn btn-sm btn-outline-danger"
                                            onClick={() => handleDelete(resource.id)}
                                            title={t('permissions.modules.delete')}
                                        >
                                            <iconify-icon icon="solar:trash-bin-minimalistic-bold"
                                                          style={{verticalAlign: 'middle'}}>
                                                <span className="ms-1">{t('permissions.modules.delete')}</span>
                                            </iconify-icon>

                                        </button>
                                    </td>
                                </tr>
                            )
                        })
                    )}
                    </tbody>
                </table>
            </div>

            {/* Modal for Create/Edit */}
            {showModal && (
                <div className="modal show d-block" style={{backgroundColor: 'rgba(0,0,0,0.5)'}}>
                    <div className="modal-dialog">
                        <div className="modal-content">
                            <div className="modal-header">
                                <h5 className="modal-title">
                                    {editingResource ? t('permissions.modules.editResource') : t('permissions.modules.addResource')}
                                </h5>
                                <button
                                    type="button"
                                    className="btn-close"
                                    onClick={handleCloseModal}
                                ></button>
                            </div>
                            <form onSubmit={handleSubmit}>
                                <div className="modal-body">
                                    <div className="mb-3">
                                        <label className="form-label">{t('permissions.modules.name')}</label>
                                        <input
                                            type="text"
                                            className="form-control"
                                            value={formData.name}
                                            onChange={(e) => setFormData({...formData, name: e.target.value})}
                                            required
                                        />
                                    </div>
                                    <div className="mb-3">
                                        <label className="form-label">{t('permissions.modules.description')}</label>
                                        <textarea
                                            className="form-control"
                                            rows="3"
                                            value={formData.description}
                                            onChange={(e) => setFormData({...formData, description: e.target.value})}
                                        />
                                    </div>
                                </div>
                                <div className="modal-footer">
                                    <button
                                        type="button"
                                        className="btn btn-secondary"
                                        onClick={handleCloseModal}
                                    >
                                        {t('permissions.modules.cancel')}
                                    </button>
                                    <button type="submit" className="btn btn-primary">
                                        {editingResource ? t('permissions.modules.update') : t('permissions.modules.create')}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

export default Modules
