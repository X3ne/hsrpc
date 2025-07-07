import { invoke } from '@tauri-apps/api/core'
import { ArrowRight2, SearchNormal } from 'iconsax-reactjs'
import { useState } from 'react'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'

import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'

import packageJson from '@/../package.json'
import { CardCta, CardCtaGroup } from '@/components/card-cta'
import { Button } from '@/components/ui/button'
import { ChangelogModal } from '@/components/modals/changelog/modal'
import { NewUpdateModal, Update } from '@/components/modals/update/modal'

const AboutSettings = () => {
  const [checkingForUpdates, setCheckingForUpdates] = useState(false)
  const [isChangelogOpen, setIsChangelogOpen] = useState(false)
  const [update, setUpdate] = useState<Update | null>(null)
  const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false)

  const version = packageJson.version
  const appName = packageJson.name

  const allDependencies = {
    ...(packageJson.dependencies || {}),
    ...(packageJson.devDependencies || {})
  }

  const dependencyList = Object.entries(allDependencies).map(([name, ver]) => ({
    name,
    version: String(ver)
  }))

  dependencyList.sort((a, b) => a.name.localeCompare(b.name))

  const checkForUpdates = async () => {
    try {
      setCheckingForUpdates(true)
      const update = await invoke('check_for_updates')
      if (update) {
        setUpdate(update as Update)
        setIsUpdateModalOpen(true)
      } else {
        setUpdate(null)
        setIsUpdateModalOpen(true)
      }
    } catch (error) {
      console.error('Error checking for updates:', error)
      toast.error('Failed to check for updates. Please try again later.')
    } finally {
      setCheckingForUpdates(false)
    }
  }

  return (
    <div className='p-6'>
      <ChangelogModal
        changelogUrl='https://raw.githubusercontent.com/X3ne/hsrpc/refs/heads/main/CHANGELOG.md'
        open={isChangelogOpen}
        onOpenChange={setIsChangelogOpen}
      />
      <NewUpdateModal
        update={update}
        open={isUpdateModalOpen}
        onOpenChange={setIsUpdateModalOpen}
      />
      <h3 className='text-xl mb-4'>About</h3>
      <div className='flex flex-col space-y-6'>
        <div className='flex flex-row items-center'>
          <img src='/icon.png' alt='Logo' className='w-10 h-10 rounded-xl' />
          <div className='ml-4 flex items-center space-x-2 h-4'>
            <p>{appName}</p>
            <Separator orientation='vertical' />
            <p>v{version}</p>
          </div>
        </div>
        <div className='flex flex-col space-y-2'>
          <p>MIT License</p>
          <p>For more information, visit our GitHub repository.</p>
        </div>

        <div>
          <h4 className='text-md text-muted-foreground mb-2'>Updates</h4>
          <CardCtaGroup>
            <CardCta
              title='View Changelog'
              actionComponent={<ArrowRight2 size={14} />}
              onClick={() => setIsChangelogOpen(true)}
            />
            <CardCta
              title='Check for Updates'
              position='bottom'
              actionComponent={
                <Button
                  variant='outline'
                  size='sm'
                  className='w-full h-full flex items-center justify-center text-[12px] px-0.5 py-1'
                  onClick={() => checkForUpdates()}
                >
                  {checkingForUpdates ? <Loader2 className='animate-spin' /> : <SearchNormal />}
                  Check Now
                </Button>
              }
            />
          </CardCtaGroup>
        </div>

        <div>
          <h4 className='text-md text-muted-foreground mb-2'>Client Log</h4>
          <CardCtaGroup>
            <CardCta
              title='View Client Log'
              actionComponent={<ArrowRight2 size={14} />}
              onClick={() => invoke('open_log_file')}
            />
          </CardCtaGroup>
        </div>

        <Separator />
        <div className='flex-1 flex flex-col min-h-0'>
          <h4 className='text-md text-muted-foreground mb-2'>Packages Used:</h4>

          <ScrollArea className='flex-1 border rounded-md p-3'>
            <ul className='list-disc list-inside space-y-1 text-sm'>
              {dependencyList.length > 0 ? (
                dependencyList.map(dep => (
                  <li key={dep.name} className='flex justify-between items-center pr-2'>
                    <span>{dep.name}</span>
                    <span className='text-xs ml-2'>{dep.version}</span>
                  </li>
                ))
              ) : (
                <p>No packages found.</p>
              )}
            </ul>
          </ScrollArea>
        </div>
      </div>
    </div>
  )
}

export { AboutSettings }
