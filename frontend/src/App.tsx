import { useState, useEffect } from 'react'
import { Onboarding } from './components/Onboarding'
import { Dashboard } from './components/Dashboard'
import './index.css'

function App() {
  const [hasOnboarded, setHasOnboarded] = useState<boolean | null>(null)

  useEffect(() => {
    // Check if user has already onboarded
    const onboarded = localStorage.getItem('hasOnboarded')
    setHasOnboarded(onboarded === 'true')
  }, [])

  // Show nothing while reading local storage to prevent flicker
  if (hasOnboarded === null) {
    return null
  }

  if (!hasOnboarded) {
    return <Onboarding onComplete={() => setHasOnboarded(true)} />
  }

  return <Dashboard />
}

export default App
