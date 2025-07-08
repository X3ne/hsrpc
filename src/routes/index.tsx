import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { listen } from '@tauri-apps/api/event'

export const Route = createFileRoute('/')({
  component: Index
})

type AppState = {
  location: State
  character: State
  menu: State
  combat: State
}

type State = {
  location: Data
  character: Data
  menu: Data
  combat: Combat
}

type Combat = {
  started: string
  is_boss: boolean
  boss: Data
}

type Data = {
  asset_id: string
  value: string
  message: string
  region: string
  sub_region: string
}

function Index() {
  const [appState, setAppState] = useState<AppState | null>(null)

  useEffect(() => {
    const unlisten = listen<AppState>('app-state', event => {
      console.log('app-state event received', event)
      setAppState(event.payload)
    })

    return () => {
      unlisten.then(f => f())
    }
  }, [])

  return (
    <div className='flex items-center justify-center h-screen overflow-hidden'>
      {/* <pre>{JSON.stringify(appState?.character, null, 2)}</pre> */}
      {appState && appState.location && (
        <div className='flex bg-secondary rounded-md overflow-hidden p-4'>
          <div className='relative pb-2 pr-2'>
            <img
              className='w-24 aspect-square rounded-md'
              src={'/assets/locations/${appState.location.location.asset_id}.png'}
            />
            <img
              className='absolute right-0 bottom-0 w-8 rounded-full'
              src={'/assets/characters/${appState.location.character.asset_id}.png'}
            />
          </div>
          <div>
            <h1 className='text-lg font-bold'>Honkai Star Rail</h1>
            <p>{appState.location.location.region}</p>
            <p>{appState.location.location.value}</p>
            {/* <p>
                {
                  // duration since appState.app_started
                  new Date(Date.now() - new Date(appState.app_started).getTime())
                    .toISOString()
                    .substr(11, 8)
                }
              </p> */}
          </div>
        </div>
      )}
    </div>
  )
}
