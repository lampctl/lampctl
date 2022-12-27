import { useEffect, useState } from 'react'
import { sortByProp } from '../lib/util'
import Group from './Group'
import styles from './Provider.module.css'

export default function Provider({ provider }) {

  const [response, setResponse] = useState(null)

  useEffect(() => {
    fetch(`/api/providers/${provider.id}`)
      .then(r => r.json())
      .then(r => setResponse(r))
  }, [])

  function setState(state) {
    fetch(`/api/providers/${provider.id}/apply/all`, {
      method: 'POST',
      headers: { 'Content-type': 'application/json' },
      body: JSON.stringify({
        state,
        duration: 1000,
      })
    })
  }

  return (
    <div className={styles.provider}>
      <div className={styles.header}>
        <div className={styles.title}>{provider.name}</div>
        <div className={styles.button} onClick={() => setState(true)}>On</div>
        <div className={styles.button} onClick={() => setState(false)}>Off</div>
      </div>
      {response === null ?
        <p>Loading...</p> :
        sortByProp(response.groups, 'name').map(g => (
          <Group
            key={g.id}
            provider={provider}
            group={g}
            lamps={response.lamps}
          />
        ))
      }
    </div>
  )
}
