import { CardCta, CardCtaGroup } from '@/components/card-cta'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Switch } from '@/components/ui/switch'

const GeneralSettings = () => {
  return (
    <>
      <div className='h-fit p-6'>
        <h3 className='text-xl'>General</h3>
      </div>
      <ScrollArea className='h-full w-full overflow-hidden'>
        <div className='flex flex-col space-y-6 p-6'>
          <CardCtaGroup>
            <CardCta
              title='About'
              description='View application information and check for updates.'
              actionComponent={<Switch />}
            />
          </CardCtaGroup>
        </div>
      </ScrollArea>
    </>
  )
}

export { GeneralSettings }
