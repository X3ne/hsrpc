import React, { useCallback } from 'react'

import { Config } from '@/providers/config-provider'

/**
 * A custom hook to manage individual fields within the Config object
 * Simplifies connecting form inputs (Switch, RadioGroup, text inputs)
 * to the Config state
 *
 * @param config The current Config object
 * @param onConfigChange The state setter from the parent (e.g., setConfig)
 * @param fieldName The key of the field in the Config object to manage
 * @param options Optional settings for type conversion (e.g., 'number' for string to number parsing)
 * @param callback Optional callback to run after the field value changes
 * @returns An object containing the field's `value` and an `onChange` handler
 */
function useConfigField<K extends keyof Config>(
  config: Config,
  onConfigChange: React.Dispatch<React.SetStateAction<Config | null>>,
  fieldName: K,
  options?: { type?: 'number' },
  callback: () => void = () => {}
) {
  const value = config[fieldName]

  const handleChange = useCallback(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (newValue: any) => {
      onConfigChange(prevConfig => {
        if (!prevConfig) return null

        let finalValue: Config[K] = newValue

        if (options?.type === 'number') {
          const parsed = typeof newValue === 'string' ? parseInt(newValue, 10) : newValue
          finalValue = (isNaN(parsed) ? 0 : parsed) as Config[K]
        }

        callback()

        return {
          ...prevConfig,
          [fieldName]: finalValue
        }
      })
    },
    [fieldName, onConfigChange, options?.type, callback]
  )

  return {
    value,
    onChange: handleChange
  }
}

export default useConfigField
