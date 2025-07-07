import { useState } from 'react'
import { Setting2 } from 'iconsax-reactjs'
import { DialogTitle } from '@radix-ui/react-dialog'

import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { SettingsSidebar } from '@/components/modals/settings/sidebar'
import { GeneralSettings } from '@/components/modals/settings/pages/general'
import { AboutSettings } from '@/components/modals/settings/pages/about'
import { ScrollArea } from '@/components/ui/scroll-area'

export type SettingsPageKey = 'general' | 'about'

const settingsPageComponents: Record<SettingsPageKey, React.ComponentType> = {
  general: GeneralSettings,
  about: AboutSettings
}

const SettingsModal = () => {
  const [activePage, setActivePage] = useState<SettingsPageKey>('general')

  const renderActivePageContent = () => {
    const ComponentToRender = settingsPageComponents[activePage]
    return <ComponentToRender />
  }

  return (
    <Dialog>
      <DialogTrigger>
        <Button variant='outline' size='icon'>
          <Setting2 size={18} />
        </Button>
      </DialogTrigger>
      <DialogContent className='flex p-0 h-[554px] w-[1000px] overflow-hidden'>
        <DialogTitle className='sr-only'>Settings</DialogTitle>
        <SidebarProvider className='h-[554px] min-h-0'>
          <SettingsSidebar activePage={activePage} onPageChange={setActivePage} />
          <ScrollArea className='w-full h-full flex-1'>
            <SidebarInset>{renderActivePageContent()}</SidebarInset>
          </ScrollArea>
        </SidebarProvider>
      </DialogContent>
    </Dialog>
  )
}

export { SettingsModal }
