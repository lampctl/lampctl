import { Outlet } from 'react-router-dom'
import Header from './Header'
import styles from './App.module.css'

export default function App() {
  return (
    <>
      <Header />
      <div className={styles.app}>
        <Outlet />
      </div>
    </>
  )
}
