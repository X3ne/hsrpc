// SettingsSidebar.tsx
import * as React from 'react'
import { Setting2, Information, type Icon } from 'iconsax-reactjs'

import { NavMain } from '@/components/modals/settings/nav'
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

export function SettingsSidebar({ activePage, onPageChange, ...props }: SettingsSidebarProps) {
  return (
    <Sidebar collapsible='icon' {...props}>
      <SidebarHeader className='px-4 mt-4'>
        <h2 className='text-xl font-semibold'>Settings</h2>
      </SidebarHeader>
      <SidebarContent>
        <NavMain
          items={sidebarNavItems.navMain}
          activePage={activePage}
          onPageChange={onPageChange}
        />
      </SidebarContent>
    </Sidebar>
  )
}
