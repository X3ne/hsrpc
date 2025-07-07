import React from 'react'
import { SideBar } from './sidebar'
import { TopBar } from './topbar'

type LayoutProps = {
  children: React.ReactNode
}

export function Layout({ children }: LayoutProps) {
  return (
    <div className='flex min-h-screen flex-col'>
      <TopBar />
      <main className='flex w-full flex-1 bg-background overflow-hidden'>
        <SideBar />
        <div className='mx-auto flex w-full'>{children}</div>
      </main>
    </div>
  )
}
