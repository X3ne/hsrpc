import { useState } from 'react'

import { CardCta, CardCtaGroup } from '@/components/card-cta'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Slider } from '@/components/ui/slider'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'

const GameSettings = () => {
  const [statusGrabbingEnabled, setStatusGrabbingEnabled] = useState(true)
  const [loopTime, setLoopTime] = useState(2000)
  const [accountUid, setAccountUid] = useState('')
  const [accountName, setAccountName] = useState('')
  const [displayAccountName, setDisplayAccountName] = useState(true)
  const [displayAccountLevel, setDisplayAccountLevel] = useState(true)
  const [preprocessThreshold, setPreprocessThreshold] = useState(135)
  const [windowName, setWindowName] = useState('Star Rail')
  const [windowClass, setWindowClass] = useState('UnityWndClass')

  return (
    <>
      <div className='h-fit p-6'>
        <h3 className='text-xl'>Game</h3>
      </div>
      <ScrollArea className='h-full w-full overflow-hidden'>
        <div className='flex flex-col space-y-6 px-6 pb-6'>
          <div>
            <CardCtaGroup>
              <CardCta
                title='Enable Status Grabbing'
                description='If you want to disable the status grabbing feature, you can do so here. This will stop the app from fetching your game status.'
                actionComponent={
                  <Switch
                    checked={statusGrabbingEnabled}
                    onCheckedChange={setStatusGrabbingEnabled}
                  />
                }
              />
              <CardCta
                title='Loop Time'
                description='Set the time interval for update loops. If you experience performance issues, try decreasing this value.'
                content={
                  <div className='flex items-center gap-2'>
                    <Slider
                      defaultValue={[2000]}
                      max={20000}
                      min={100}
                      step={100}
                      value={[loopTime]}
                      onValueChange={value => setLoopTime(value[0])}
                      className='w-full'
                    />
                    <Input
                      type='number'
                      value={loopTime}
                      onChange={e => setLoopTime(Number(e.target.value))}
                      className='w-20'
                      min={100}
                      max={20000}
                      step={100}
                      aria-label='Loop Time'
                    />
                  </div>
                }
              />
            </CardCtaGroup>
          </div>
          <div>
            <h4 className='text-md text-muted-foreground mb-2'>Hsr Account</h4>

            <CardCtaGroup>
              <CardCta
                title='Account UID'
                description='Set your hsr account UID to enable features like displaying your in-game level.'
                content={
                  <Input
                    type='number'
                    placeholder='700008442'
                    aria-label='Account UID'
                    value={accountUid}
                    onChange={e => setAccountUid(e.target.value)}
                  />
                }
              />
              <CardCta
                title='Account Name'
                description='Set your hsr account name to enable trailblazer detection. If you have set your account UID, you can leave this blank.'
                content={
                  <Input
                    type='text'
                    placeholder='Account Name'
                    aria-label='Account Name'
                    value={accountName}
                    onChange={e => setAccountName(e.target.value)}
                  />
                }
              />
              <CardCta
                title='Display Account Name'
                description='Enable this to show your account name in the Discord presence.'
                actionComponent={
                  <Switch checked={displayAccountName} onCheckedChange={setDisplayAccountName} />
                }
              />
              <CardCta
                title='Display Account Level'
                description='Enable this to show your account level in the Discord presence.'
                actionComponent={
                  <Switch checked={displayAccountLevel} onCheckedChange={setDisplayAccountLevel} />
                }
              />
            </CardCtaGroup>
          </div>
          <div>
            <h4 className='text-md text-muted-foreground mb-2'>Advanced</h4>

            <CardCta
              title='Preprocess Threshold'
              description='Set the threshold for preprocessing ocr images. If you experience issues with the ocr, try tweaking this value.'
              content={
                <div className='flex items-center gap-2'>
                  <Slider
                    defaultValue={[135]}
                    max={1000}
                    min={0}
                    value={[preprocessThreshold]}
                    onValueChange={value => setPreprocessThreshold(value[0])}
                    className='w-full'
                  />
                  <Input
                    type='number'
                    value={preprocessThreshold}
                    onChange={e => setPreprocessThreshold(Number(e.target.value))}
                    className='w-20'
                    min={0}
                    max={1000}
                    aria-label='Preprocess Threshold'
                  />
                </div>
              }
            />

            <CardCtaGroup>
              <CardCta
                title='Window Name'
                description='Set the window name for the game. This is used to identify the game window. Check the README on the GitHub repository for more information.'
                content={
                  <Input
                    type='text'
                    placeholder='Star Rail'
                    defaultValue={'Star Rail'}
                    aria-label='Window Name'
                    value={windowName}
                    onChange={e => setWindowName(e.target.value)}
                  />
                }
              />
              <CardCta
                title='Window Class'
                description='Set the window class for the game. This is used to identify the game window. Check the README on the GitHub repository for more information.'
                content={
                  <Input
                    type='text'
                    placeholder='UnityWndClass'
                    defaultValue={'UnityWndClass'}
                    aria-label='Window Class'
                    value={windowClass}
                    onChange={e => setWindowClass(e.target.value)}
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

export { GameSettings }
