import { useState } from 'react'
import { Setting2 } from 'iconsax-reactjs'
import { DialogTitle } from '@radix-ui/react-dialog'

import { Dialog, DialogContent, DialogTrigger } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { SettingsSidebar } from '@/components/modals/settings/sidebar'
import { GeneralSettings } from '@/components/modals/settings/pages/general'
import { AboutSettings } from '@/components/modals/settings/pages/about'

type SettingsPageKey = 'general' | 'about'

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
      <DialogContent className='p-0 h-9/12 w-3/4'>
        <DialogTitle className='sr-only'>Settings</DialogTitle>
        <SidebarProvider className='min-h-0'>
          <SettingsSidebar
            activePage={activePage}
            onPageChange={setActivePage}
            className='h-full'
          />
          <SidebarInset className='h-full rounded-lg'>{renderActivePageContent()}</SidebarInset>
        </SidebarProvider>
      </DialogContent>
    </Dialog>
  )
}

export { SettingsModal, type SettingsPageKey }
