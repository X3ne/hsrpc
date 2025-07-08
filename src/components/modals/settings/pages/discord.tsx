import { useState } from 'react'

import { CardCta, CardCtaGroup } from '@/components/card-cta'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-react'
import { Refresh } from 'iconsax-reactjs'
import { Config } from '@/types'
import useConfigField from '@/hooks/use-config-fields'

interface DiscordSettingsProps {
  config: Config
  onConfigChange: React.Dispatch<React.SetStateAction<Config | null>>
}

const DiscordSettings: React.FC<DiscordSettingsProps> = ({ config, onConfigChange }) => {
  const [isReconnecting, setIsReconnecting] = useState(false)

  const { value: discordAppId, onChange: handleDiscordAppIdChange } = useConfigField(
    config,
    onConfigChange,
    'discord_app_id'
  )

  const handleDiscordReconnect = () => {
    setIsReconnecting(true)
    // TODO: Implement actual reconnection logic
    setTimeout(() => {
      setIsReconnecting(false)
    }, 2000)
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
                className='w-full h-full flex items-center justify-center text-[12px] px-0.5 py-1'
                onClick={() => handleDiscordReconnect()}
              >
                {isReconnecting ? <Loader2 className='animate-spin' /> : <Refresh />}
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
