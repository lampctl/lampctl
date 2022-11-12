import styles from './Group.module.css'
import Lamp from './Lamp'

export default function Group({ provider, group, lamps }) {
  return (
    <div className={styles.group}>
      <div className={styles.title}>{group.name}</div>
      <div className={styles.lamps}>
        {lamps.filter(l => l.group_id == group.id).map(l => (
          <Lamp
            key={l.id}
            provider={provider}
            group={group}
            lamp={l}
          />
        ))}
      </div>
    </div>
  )
}
