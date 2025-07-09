import { Link } from '@tanstack/react-router'
import { SettingsModal } from '@/components/modals/settings/modal'

const SideBar = () => {
  return (
    <div className='flex flex-col items-center justify-between h-screen w-16 bg-card p-3'>
      <Link to='/' className='z-50'>
        <img src='/icon.png' className='w-8 h-8 mt-2 aspect-square' />
      </Link>

      <SettingsModal />
    </div>
  )
}

export { SideBar }
