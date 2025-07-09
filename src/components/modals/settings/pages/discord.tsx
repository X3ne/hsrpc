import React from 'react'

import { CardCta, CardCtaGroup } from '@/components/card-cta'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-react'
import { Refresh } from 'iconsax-reactjs'
import { Config } from '@/providers/config-provider'
import useConfigField from '@/hooks/use-config-fields'
import { invoke } from '@tauri-apps/api/core'
import { cn } from '@/lib/utils'

interface DiscordSettingsProps {
  config: Config
  onConfigChange: React.Dispatch<React.SetStateAction<Config | null>>
}

type ReconnectionStatus = 'reconnecting' | 'error' | 'connected'

const DiscordSettings: React.FC<DiscordSettingsProps> = ({ config, onConfigChange }) => {
  const [reconnectionStatus, setReconnectionStatus] =
    React.useState<ReconnectionStatus>('connected')

  const { value: discordAppId, onChange: handleDiscordAppIdChange } = useConfigField(
    config,
    onConfigChange,
    'discord_app_id'
  )

  React.useEffect(() => {
    const checkDiscordConnection = async () => {
      try {
        const isConnected = await invoke<boolean>('is_ipc_connected')
        if (isConnected) {
          setReconnectionStatus('connected')
        } else {
          setReconnectionStatus('error')
        }
      } catch (e) {
        console.error('Failed to check Discord connection:', e)
        setReconnectionStatus('error')
      }
    }
    checkDiscordConnection()
  }, [])

  const handleDiscordReconnect = async () => {
    try {
      setReconnectionStatus('reconnecting')
      await invoke('reconnect_to_discord')
      setTimeout(() => {
        setReconnectionStatus('connected')
      }, 1000)
    } catch (e) {
      console.error('Failed to reconnect to Discord:', e)
      setReconnectionStatus('error')
    }
  }

  const renderReconnectionStatusIcon = (status: ReconnectionStatus) => {
    switch (status) {
      case 'connected':
        return <Refresh />
      case 'reconnecting':
        return <Loader2 className='animate-spin' />
      case 'error':
        return <Refresh />
    }
  }

  return (
    <>
      <div className='h-fit p-6'>
        <h3 className='text-xl'>Discord</h3>
      </div>
      <ScrollArea className='h-full w-full overflow-hidden'>
        <div className='flex flex-col space-y-6 px-6 pb-6'>
          <CardCta
            title='Reconnect to Discord'
            actionComponent={
              <Button
                variant='outline'
                size='sm'
                className={cn(
                  'w-full h-full flex items-center justify-center text-[12px] px-0.5 py-1',
                  reconnectionStatus === 'error' &&
                    '!border-destructive !text-destructive-foreground !bg-destructive/20 hover:!bg-destructive/30'
                )}
                disabled={reconnectionStatus === 'reconnecting'}
                onClick={() => handleDiscordReconnect()}
              >
                {renderReconnectionStatusIcon(reconnectionStatus)}
                Reconnect
              </Button>
            }
          />

          <div>
            <h4 className='text-md text-muted-foreground mb-2'>Advanced</h4>

            <CardCtaGroup>
              <CardCta
                title='Discord App ID'
                description='You can use your own Discord App ID to personalize assets or testing.'
                content={
                  <Input
                    type='number'
                    placeholder='1208212792574869544'
                    value={discordAppId}
                    onChange={e => handleDiscordAppIdChange(e.target.value || null)}
                    aria-label='Discord App ID'
                  />
                }
              />
            </CardCtaGroup>
          </div>
        </div>
      </ScrollArea>
    </>
  )
}

export { DiscordSettings }
