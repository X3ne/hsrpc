import * as React from 'react'

import packageJson from '@/../package.json'

export function useVersion() {
  const [version, setVersion] = React.useState<string>(packageJson.version)

  React.useEffect(() => {
    setVersion(packageJson.version)
  }, [])

  return { version }
}
