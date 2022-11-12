import { useState } from 'react'
import Toggle from './Toggle'
import styles from './Lamp.module.css'

export default function Lamp({ provider, group, lamp }) {

  const [state, setState] = useState(lamp.state)

  function handleChange() {
    fetch(`/api/providers/${provider.id}/apply`, {
      method: 'POST',
      headers: { 'Content-type': 'application/json' },
      body: JSON.stringify([{
        group_id: group.id,
        lamp_id: lamp.id,
        state: !state
      }])
    })
      .then(() => {
        setState(!state)
      })
  }

  return (
    <div className={styles.container}>
      <div className={styles.lamp}>
        {lamp.name}
        <Toggle state={state} onChange={handleChange} />
      </div>
    </div>
  )
}
