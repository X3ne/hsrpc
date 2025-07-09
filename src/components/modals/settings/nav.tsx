import { type Icon } from 'iconsax-reactjs'

import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem
} from '@/components/ui/sidebar'
import { SettingsPageKey } from './modal'

interface SettingsModalNavProps {
  items: {
    title: string
    key: SettingsPageKey
    icon?: Icon
  }[]
  activePage: SettingsPageKey
  onPageChange: (pageKey: SettingsPageKey) => void
}

const SettingsModalNav = ({ items, activePage, onPageChange }: SettingsModalNavProps) => {
  return (
    <SidebarGroup className='py-0'>
      <SidebarMenu>
        {items.map(item => (
          <SidebarMenuItem key={item.key}>
            <SidebarMenuButton
              tooltip={item.title}
              isActive={item.key === activePage}
              onClick={() => onPageChange(item.key)}
            >
              {item.icon && <item.icon className='scale-120' />}
              <span>{item.title}</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  )
}

export { SettingsModalNav, type SettingsModalNavProps }
