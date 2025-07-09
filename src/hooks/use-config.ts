import { ConfigContext } from '@/providers/config-provider'
import { useContext } from 'react'

export const useConfig = () => {
  const context = useContext(ConfigContext)

  if (!context) {
    throw new Error('useConfig must be used within a ConfigProvider')
  }

  return context
}
