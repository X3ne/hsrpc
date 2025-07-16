import { UpdateContext } from '@/providers/update-provider'
import { useContext } from 'react'

export const useUpdate = () => {
  const context = useContext(UpdateContext)

  if (!context) {
    throw new Error('useUpdate must be used within a UpdateContext')
  }

  return context
}
