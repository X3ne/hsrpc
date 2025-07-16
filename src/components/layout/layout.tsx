import React from 'react'

import { Toaster } from '@/components/ui/sonner'

// import { SideBar } from './sidebar'
import { TopBar } from './topbar'
import { NewUpdateModal } from '../modals/update/modal'

type LayoutProps = {
  children: React.ReactNode
}

const Layout = ({ children }: LayoutProps) => {
  return (
    <div className='flex min-h-screen flex-col'>
      <TopBar />
      <main className='flex w-full flex-1 bg-background overflow-hidden'>
        {/* <SideBar /> */}
        <div className='mx-auto flex w-full'>{children}</div>
      </main>
      <NewUpdateModal />
      <Toaster />
    </div>
  )
}

export { Layout, type LayoutProps }
