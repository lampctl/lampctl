import styles from './Toggle.module.css'

export default function Toggle({ state, onChange }) {

  let classNames = [styles.slider]

  if (state) {
    classNames.push(styles.on)
  }

  return (
    <div className={styles.toggle} onClick={onChange}>
      <div className={classNames.join(' ')}>
        {state ? "ON" : "OFF"}
      </div>
    </div>
  )
}
