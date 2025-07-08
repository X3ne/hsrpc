import * as React from 'react'
import { Setting2, Information, Gameboy, Activity, type Icon } from 'iconsax-reactjs'

import { SettingsModalNav } from '@/components/modals/settings/nav'
import { Sidebar, SidebarContent, SidebarHeader } from '@/components/ui/sidebar'
import { SettingsPageKey } from './modal'

export type ISidebarNavItems = {
  navMain: {
    title: string
    key: SettingsPageKey
    icon?: Icon
  }[]
}

const sidebarNavItems: ISidebarNavItems = {
  navMain: [
    {
      title: 'General',
      key: 'general',
      icon: Setting2
    },
    {
      title: 'Game',
      key: 'game',
      icon: Gameboy
    },
    {
      title: 'Discord',
      key: 'discord',
      icon: Activity
    },
    {
      title: 'About',
      key: 'about',
      icon: Information
    }
  ]
}

interface SettingsSidebarProps extends React.ComponentProps<typeof Sidebar> {
  activePage: SettingsPageKey
  onPageChange: (pageKey: SettingsPageKey) => void
}

const SettingsSidebar = ({ activePage, onPageChange, ...props }: SettingsSidebarProps) => {
  return (
    <Sidebar {...props}>
      <SidebarHeader className='p-0'>
        <h2 className='text-xl px-6 py-5 text-muted-foreground'>Settings</h2>
      </SidebarHeader>
      <SidebarContent>
        <SettingsModalNav
          items={sidebarNavItems.navMain}
          activePage={activePage}
          onPageChange={onPageChange}
        />
      </SidebarContent>
    </Sidebar>
  )
}

export { SettingsSidebar, type SettingsSidebarProps }
