import { Button } from '@/components/ui/button'
import { getCurrentWindow } from '@tauri-apps/api/window'
import { Minus, X } from 'lucide-react'

const TopBar = () => {
  return (
    <div
      data-tauri-drag-region
      className='absolute top-0 left-0 flex items-center justify-end w-full p-2 gap-4 z-50'
    >
      <Button
        variant='ghost'
        className='w-6 h-6 rounded-[6px]'
        size='icon'
        onClick={() => getCurrentWindow().minimize()}
      >
        <Minus color='white' className='w-4' />
      </Button>
      <Button
        variant='ghost'
        className='w-6 h-6 rounded-[6px]'
        size='icon'
        onClick={() => getCurrentWindow().hide()}
      >
        <X color='white' className='w-4' />
      </Button>
    </div>
  )
}

export { TopBar }
