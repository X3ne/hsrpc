import { type Icon } from 'iconsax-reactjs'

import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem
} from '@/components/ui/sidebar'
import { SettingsPageKey } from './modal'

interface NavMainProps {
  items: {
    title: string
    key: SettingsPageKey
    icon?: Icon
  }[]
  activePage: SettingsPageKey
  onPageChange: (pageKey: SettingsPageKey) => void
}

export function NavMain({ items, activePage, onPageChange }: NavMainProps) {
  return (
    <SidebarGroup>
      <SidebarMenu>
        {items.map(item => (
          <SidebarMenuItem key={item.key}>
            <SidebarMenuButton
              tooltip={item.title}
              isActive={item.key === activePage}
              onClick={() => onPageChange(item.key)}
            >
              {item.icon && <item.icon size={18} />}
              <span>{item.title}</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  )
}
