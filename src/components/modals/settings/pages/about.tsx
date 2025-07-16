import { invoke } from '@tauri-apps/api/core'
import { openUrl } from '@tauri-apps/plugin-opener'
import { ArrowRight2 } from 'iconsax-reactjs'
import { useState } from 'react'
import { Loader2, Search, SquareArrowOutUpRight } from 'lucide-react'
import { toast } from 'sonner'
import { error } from '@tauri-apps/plugin-log'

import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { CardCta, CardCtaGroup } from '@/components/card-cta'
import { Button } from '@/components/ui/button'
import { ChangelogModal } from '@/components/modals/changelog/modal'
import packageJson from '@/../package.json'
import { Update } from '@/providers/update-provider'
import { useUpdate } from '@/hooks/use-update'

const AboutSettings = () => {
  const [checkingForUpdates, setCheckingForUpdates] = useState(false)
  const [isChangelogOpen, setIsChangelogOpen] = useState(false)
  const { setUpdate } = useUpdate()

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
      } else {
        setUpdate(null)
      }
    } catch (e) {
      error('Error checking for updates:', {
        keyValues: {
          error: e instanceof Error ? e.message : String(e)
        }
      })
      toast.error('Failed to check for updates. Please check logs and report this issue.')
    } finally {
      setCheckingForUpdates(false)
    }
  }

  return (
    <>
      <ChangelogModal
        changelogUrl='https://raw.githubusercontent.com/X3ne/hsrpc/refs/heads/main/CHANGELOG.md'
        open={isChangelogOpen}
        onOpenChange={setIsChangelogOpen}
      />
      <div className='h-fit p-6'>
        <h3 className='text-xl'>About</h3>
      </div>
      <ScrollArea className='h-full w-full overflow-hidden'>
        <div className='flex flex-col space-y-6 px-6 pb-6'>
          <div className='flex flex-row items-center'>
            <img src='/icon.png' alt='Logo' className='w-10 h-10 rounded-xl' />
            <div className='ml-4 flex items-center space-x-2 h-4'>
              <p>{appName}</p>
              <Separator orientation='vertical' />
              <p>v{version}</p>
            </div>
          </div>
          <div className='w-fit'>
            <Button
              variant={'link'}
              className='!p-0 text-sm'
              onClick={async () => await openUrl('https://github.com/X3ne/hsrpc/blob/main/LICENSE')}
            >
              <p>AGPL-3.0 License</p>
              <SquareArrowOutUpRight size={14} />
            </Button>
            <p className='text-sm'>
              For more information, visit our{' '}
              <Button
                variant={'link'}
                className='!p-0 text-sm'
                onClick={() => openUrl('https://github.com/X3ne/hsrpc')}
              >
                <span>GitHub repository.</span>
                <SquareArrowOutUpRight size={14} />
              </Button>
            </p>
            <p className='text-sm'>
              If you like this app and you can, please consider supporting me.{' '}
              <Button
                variant={'link'}
                className='!p-0 text-sm'
                onClick={() => openUrl('https://ko-fi.com/ncrl_')}
              >
                <span>Ko-Fi</span>
                <SquareArrowOutUpRight size={14} />
              </Button>
            </p>
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
                    {checkingForUpdates ? <Loader2 className='animate-spin' /> : <Search />}
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

            <div className='flex-1 border rounded-md p-3'>
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
            </div>
          </div>
        </div>
      </ScrollArea>
    </>
  )
}

export { AboutSettings }
