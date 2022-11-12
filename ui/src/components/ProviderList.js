import { useEffect, useState } from 'react'

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
          <div key={p.id}>{p.name}</div>
        ))}
    </>
  )
}
