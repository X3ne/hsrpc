import { enable, isEnabled, disable } from '@tauri-apps/plugin-autostart'
import { error } from '@tauri-apps/plugin-log'
import { toast } from 'sonner'

import { CardCta, CardCtaGroup } from '@/components/card-cta'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Switch } from '@/components/ui/switch'
import { ClosingBehavior, Config } from '@/providers/config-provider'
import useConfigField from '@/hooks/use-config-fields'

interface GeneralSettingsProps {
  config: Config
  onConfigChange: React.Dispatch<React.SetStateAction<Config | null>>
}

const GeneralSettings: React.FC<GeneralSettingsProps> = ({ config, onConfigChange }) => {
  const { value: autostart, onChange: handleAutostartChange } = useConfigField(
    config,
    onConfigChange,
    'autostart',
    undefined,
    async () => {
      try {
        if (await isEnabled()) {
          await disable()
        } else {
          await enable()
        }
      } catch (e) {
        error('Failed to toggle autostart:', {
          keyValues: {
            error: e instanceof Error ? e.message : String(e)
          }
        })
        toast.error('Failed to toggle autostart. Please check logs and report this issue.')
      }
    }
  )

  const { value: trayLaunch, onChange: handleTrayLaunchChange } = useConfigField(
    config,
    onConfigChange,
    'tray_launch'
  )

  const { value: closingBehavior, onChange: handleClosingBehaviorChange } = useConfigField(
    config,
    onConfigChange,
    'closing_behavior'
  )

  const { value: autoUpdate, onChange: handleAutoUpdateChange } = useConfigField(
    config,
    onConfigChange,
    'auto_update'
  )

  return (
    <>
      <div className='h-fit p-6'>
        <h3 className='text-xl'>General</h3>
      </div>
      <ScrollArea className='h-full w-full overflow-hidden'>
        <div className='flex flex-col space-y-6 px-6 pb-6'>
          <div>
            <h4 className='text-md text-muted-foreground mb-2'>Launch Settings</h4>
            <CardCtaGroup>
              <CardCta
                title='Auto Launch on Startup'
                actionComponent={
                  <Switch checked={autostart} onCheckedChange={handleAutostartChange} />
                }
              />
              <CardCta
                title={
                  <div className='flex items-center gap-2'>
                    <p>Launch in System Tray</p>
                    <Badge variant='secondary'>Recommended</Badge>
                  </div>
                }
                description='Launch the application in the system tray.'
                actionComponent={
                  <Switch checked={trayLaunch} onCheckedChange={handleTrayLaunchChange} />
                }
              />
            </CardCtaGroup>
          </div>
          <div>
            <h4 className='text-md text-muted-foreground mb-2'>Close Settings</h4>
            <CardCtaGroup>
              <CardCta
                title='Close App'
                content={
                  <RadioGroup
                    value={closingBehavior}
                    onValueChange={value => handleClosingBehaviorChange(value as ClosingBehavior)}
                  >
                    <div className='flex items-center space-x-2'>
                      <RadioGroupItem value='minimize' id='minimize' />
                      <Label htmlFor='tray'>Minimize to Tray</Label>
                      <Badge variant='secondary'>Recommended</Badge>
                    </div>
                    <div className='flex items-center space-x-2'>
                      <RadioGroupItem value='exit' id='exit' />
                      <Label htmlFor='quit'>Quit Application</Label>
                    </div>
                  </RadioGroup>
                }
              />
            </CardCtaGroup>
          </div>
          <div>
            <h4 className='text-md text-muted-foreground mb-2'>Updates</h4>
            <CardCtaGroup>
              <CardCta
                title='Auto Update App on Start'
                description='Automatically check for updates and install them on application start.'
                actionComponent={
                  <Switch checked={autoUpdate} onCheckedChange={handleAutoUpdateChange} />
                }
              />
            </CardCtaGroup>
          </div>
        </div>
      </ScrollArea>
    </>
  )
}

export { GeneralSettings }
