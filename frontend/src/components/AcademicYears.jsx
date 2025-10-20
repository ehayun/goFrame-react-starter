import React, { useState, useEffect } from 'react'
import Select from 'react-select'

const AcademicYears = () => {
  const [selectedYear, setSelectedYear] = useState(null)
  const [academicYears, setAcademicYears] = useState([])
  const [loading, setLoading] = useState(true)

  // Load academic years list and saved academic year on component mount
  useEffect(() => {
    loadAcademicYearsList()
  }, [])

  const loadAcademicYearsList = async () => {
    try {
      // First, load the list of available academic years
      const yearsResponse = await fetch('/api/academic-years', {
        method: 'GET',
        credentials: 'include'
      })
      
      if (yearsResponse.ok) {
        const yearsData = await yearsResponse.json()
        if (yearsData.academicYears) {
          // Convert array of strings to react-select options
          const yearOptions = yearsData.academicYears.map(year => ({
            value: year,
            label: year
          }))
          setAcademicYears(yearOptions)
          
          // Then load the saved academic year
          await loadAcademicYear(yearOptions)
        } else {
          setLoading(false)
        }
      } else {
        setLoading(false)
      }
    } catch (error) {
      console.error('Error loading academic years list:', error)
      setLoading(false)
    }
  }

  const loadAcademicYear = async (yearOptions = academicYears) => {
    try {
      const response = await fetch('/api/academic-year', {
        method: 'GET',
        credentials: 'include'
      })
      
      if (response.ok) {
        const data = await response.json()
        if (data.academicYear && data.academicYear !== '') {
          const yearOption = yearOptions.find(year => year.value === data.academicYear)
          if (yearOption) {
            setSelectedYear(yearOption)
          } else {
            // If saved year not found in current list, use the first year
            if (yearOptions.length > 0) {
              setSelectedYear(yearOptions[0])
            }
          }
        } else {
          // If no key in Redis or empty value, use the first year in the list
          if (yearOptions.length > 0) {
            setSelectedYear(yearOptions[0])
          }
        }
      } else {
        // If API call fails, use the first year in the list
        if (yearOptions.length > 0) {
          setSelectedYear(yearOptions[0])
        }
      }
    } catch (error) {
      console.error('Error loading academic year:', error)
      // If error occurs, use the first year in the list
      if (yearOptions.length > 0) {
        setSelectedYear(yearOptions[0])
      }
    } finally {
      setLoading(false)
    }
  }

  const handleYearChange = async (selectedOption) => {
    try {
      const response = await fetch('/api/academic-year', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({ academicYear: selectedOption.value })
      })

      if (response.ok) {
        setSelectedYear(selectedOption)
        // Reload the page after successful save
        window.location.reload()
      } else {
        console.error('Failed to save academic year')
      }
    } catch (error) {
      console.error('Error saving academic year:', error)
    }
  }

  const customStyles = {
    control: (provided, state) => ({
      ...provided,
      minHeight: '32px',
      fontSize: '14px',
      borderColor: state.isFocused ? '#3b82f6' : '#d1d5db',
      boxShadow: state.isFocused ? '0 0 0 1px #3b82f6' : 'none',
      '&:hover': {
        borderColor: '#3b82f6'
      }
    }),
    valueContainer: (provided) => ({
      ...provided,
      padding: '0 8px'
    }),
    input: (provided) => ({
      ...provided,
      margin: '0px'
    }),
    indicatorSeparator: () => ({
      display: 'none'
    }),
    indicatorsContainer: (provided) => ({
      ...provided,
      height: '32px'
    }),
    dropdownIndicator: (provided) => ({
      ...provided,
      padding: '4px 8px'
    }),
    menu: (provided) => ({
      ...provided,
      fontSize: '14px'
    }),
    option: (provided, state) => ({
      ...provided,
      backgroundColor: state.isSelected ? '#3b82f6' : state.isFocused ? '#f3f4f6' : 'white',
      color: state.isSelected ? 'white' : '#374151',
      '&:hover': {
        backgroundColor: state.isSelected ? '#3b82f6' : '#f3f4f6'
      }
    })
  }

  if (loading) {
    return (
      <div className="academic-years-loading">
        <div className="spinner-border spinner-border-sm" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    )
  }

  return (
    <div className="academic-years-container">
      <Select
        value={selectedYear}
        onChange={handleYearChange}
        options={academicYears}
        placeholder="Select Academic Year"
        isSearchable={false}
        styles={customStyles}
        className="academic-years-select"
      />
    </div>
  )
}

export default AcademicYears
