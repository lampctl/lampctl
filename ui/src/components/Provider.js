import { useEffect, useState } from 'react'
import Group from './Group'
import styles from './Provider.module.css'

export default function Provider({ provider }) {

  const [response, setResponse] = useState(null)

  useEffect(() => {
    fetch(`/api/providers/${provider.id}`)
      .then(r => r.json())
      .then(r => setResponse(r))
  }, [])

  return (
    <>
      <div className={styles.title}>{provider.name}</div>
      {response === null ?
        <p>Loading...</p> :
        response.groups.map(g => (
          <Group
            key={g.id}
            provider={provider}
            group={g}
            lamps={response.lamps}
          />
        ))
      }
    </>
  )
}
