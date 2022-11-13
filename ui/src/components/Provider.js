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

  return (
    <>
      <div className={styles.title}>{provider.name}</div>
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
    </>
  )
}
