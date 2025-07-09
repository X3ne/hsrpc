import { CardCta, CardCtaGroup } from '@/components/card-cta'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Slider } from '@/components/ui/slider'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Config } from '@/providers/config-provider'
import useConfigField from '@/hooks/use-config-fields'

interface GameSettingsProps {
  config: Config
  onConfigChange: React.Dispatch<React.SetStateAction<Config | null>>
}

const GameSettings: React.FC<GameSettingsProps> = ({ config, onConfigChange }) => {
  const { value: statusEnabled, onChange: handleStatusEnabledChange } = useConfigField(
    config,
    onConfigChange,
    'enable_status'
  )

  const { value: loopTime, onChange: handleLoopTimeChange } = useConfigField(
    config,
    onConfigChange,
    'loop_time'
  )

  const { value: accountUid, onChange: handleAccountUidChange } = useConfigField(
    config,
    onConfigChange,
    'account_uid'
  )

  const { value: accountName, onChange: handleAccountNameChange } = useConfigField(
    config,
    onConfigChange,
    'account_name'
  )

  const { value: displayAccountName, onChange: handleDisplayAccountNameChange } = useConfigField(
    config,
    onConfigChange,
    'display_name'
  )

  const { value: displayAccountLevel, onChange: handleDisplayAccountLevelChange } = useConfigField(
    config,
    onConfigChange,
    'display_level'
  )

  const { value: preprocessThreshold, onChange: handlePreprocessThresholdChange } = useConfigField(
    config,
    onConfigChange,
    'preprocess_threshold',
    { type: 'number' }
  )

  const { value: windowName, onChange: handleWindowNameChange } = useConfigField(
    config,
    onConfigChange,
    'window_name'
  )

  // const { value: windowClass, onChange: handleWindowClassChange } = useConfigField(
  //   config,
  //   onConfigChange,
  //   'window_class'
  // )

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
                  <Switch checked={statusEnabled} onCheckedChange={handleStatusEnabledChange} />
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
                      onValueChange={value => handleLoopTimeChange(value[0])}
                      className='w-full'
                    />
                    <Input
                      type='number'
                      value={loopTime}
                      onChange={e => handleLoopTimeChange(Number(e.target.value))}
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
                    value={accountUid || ''}
                    onChange={e => handleAccountUidChange(e.target.value || null)}
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
                    value={accountName || ''}
                    onChange={e => handleAccountNameChange(e.target.value || null)}
                  />
                }
              />
              <CardCta
                title='Display Account Name'
                description='Enable this to show your account name in the Discord presence.'
                actionComponent={
                  <Switch
                    checked={displayAccountName}
                    onCheckedChange={handleDisplayAccountNameChange}
                  />
                }
              />
              <CardCta
                title='Display Account Level'
                description='Enable this to show your account level in the Discord presence.'
                actionComponent={
                  <Switch
                    checked={displayAccountLevel}
                    onCheckedChange={handleDisplayAccountLevelChange}
                  />
                }
              />
            </CardCtaGroup>
          </div>
          <div>
            <h4 className='text-md text-muted-foreground mb-2'>Advanced</h4>

            <CardCtaGroup>
              <CardCta
                title='Preprocess Threshold'
                description='Set the threshold for preprocessing ocr images. If you experience issues with the ocr, try tweaking this value.'
                content={
                  <div className='flex items-center gap-2'>
                    <Slider
                      defaultValue={[135]}
                      max={255}
                      min={0}
                      value={[preprocessThreshold]}
                      onValueChange={value => handlePreprocessThresholdChange(value[0])}
                      className='w-full'
                    />
                    <Input
                      type='number'
                      value={preprocessThreshold}
                      onChange={e => handlePreprocessThresholdChange(Number(e.target.value))}
                      className='w-20'
                      min={0}
                      max={255}
                      aria-label='Preprocess Threshold'
                    />
                  </div>
                }
              />

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
                    onChange={e => handleWindowNameChange(e.target.value)}
                  />
                }
              />
              {/* <CardCta
                title='Window Class'
                description='Set the window class for the game. This is used to identify the game window. Check the README on the GitHub repository for more information.'
                content={
                  <Input
                    type='text'
                    placeholder='UnityWndClass'
                    defaultValue={'UnityWndClass'}
                    aria-label='Window Class'
                    value={windowClass}
                    onChange={e => handleWindowClassChange(e.target.value)}
                  />
                }
              /> */}
            </CardCtaGroup>
          </div>
        </div>
      </ScrollArea>
    </>
  )
}

export { GameSettings }
