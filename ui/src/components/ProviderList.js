import { useEffect, useState } from 'react'
import Provider from './Provider'

export default function ProviderList() {

  const [providers, setProviders] = useState(null)

  useEffect(() => {
    fetch('/api/providers')
      .then(r => r.json())
      .then(p => setProviders(p))
  }, [])

  return (
    <>
      {providers === null ?
        <p>Loading...</p> :
        providers.map(p => (
          <Provider key={p.id} provider={p} />
        ))}
    </>
  )
}
