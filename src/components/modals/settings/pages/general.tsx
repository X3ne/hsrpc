import { CardCta, CardCtaGroup } from '@/components/card-cta'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Switch } from '@/components/ui/switch'

const GeneralSettings = () => {
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
              <CardCta title='Auto Launch on Startup' actionComponent={<Switch />} />
              <CardCta
                title={
                  <div className='flex items-center gap-2'>
                    <p>Launch in System Tray</p>
                    <Badge variant='secondary'>Recommended</Badge>
                  </div>
                }
                description='Launch the application in the system tray'
                actionComponent={<Switch />}
              />
            </CardCtaGroup>
          </div>
          <div>
            <h4 className='text-md text-muted-foreground mb-2'>Close Settings</h4>
            <CardCtaGroup>
              <CardCta
                title='Close App'
                content={
                  <RadioGroup defaultValue='tray'>
                    <div className='flex items-center space-x-2'>
                      <RadioGroupItem value='tray' id='tray' />
                      <Label htmlFor='tray'>Minimize to Tray</Label>
                      <Badge variant='secondary'>Recommended</Badge>
                    </div>
                    <div className='flex items-center space-x-2'>
                      <RadioGroupItem value='quit' id='quit' />
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
                actionComponent={<Switch />}
              />
            </CardCtaGroup>
          </div>
        </div>
      </ScrollArea>
    </>
  )
}

export { GeneralSettings }
