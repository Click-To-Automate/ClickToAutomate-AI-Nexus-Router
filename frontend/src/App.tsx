import { useState, useEffect } from 'react'
import { Onboarding } from './components/Onboarding'
import { Dashboard } from './components/Dashboard'
import { Preloader } from './components/Preloader'
import './index.css'

function App() {
  const [hasOnboarded, setHasOnboarded] = useState<boolean | null>(null)
  const [preloaderComplete, setPreloaderComplete] = useState<boolean>(false)

  useEffect(() => {
    const onboarded = localStorage.getItem('hasOnboarded')
    setHasOnboarded(onboarded === 'true')
  }, [])

  if (hasOnboarded === null) {
    return null
  }

  return (
    <>
      <Preloader onComplete={() => setPreloaderComplete(true)} />
      
      {/* Only show the rest of the app if preloader is complete to avoid flashes underneath, or we can just render them below since preloader has absolute z-index. Let's render them always so they are ready behind the curtain. */}
      {preloaderComplete && (
        !hasOnboarded ? (
          <Onboarding onComplete={() => setHasOnboarded(true)} />
        ) : (
          <Dashboard />
        )
      )}
    </>
  )
}

export default App
