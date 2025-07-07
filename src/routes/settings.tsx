import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/settings')({
  component: () => (
    <div className='px-4 py-12'>
      <div className='flex items-center space-x-2'>
        <Switch id='airplane-mode' />
        <Label htmlFor='airplane-mode'>Launch at startup</Label>
      </div>
    </div>
  )
})
