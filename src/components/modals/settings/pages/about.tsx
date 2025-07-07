import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'

import packageJson from '@/../package.json'

export function AboutSettings() {
  const version = packageJson.version
  const appName = packageJson.name

  const allDependencies = {
    ...(packageJson.dependencies || {}),
    ...(packageJson.devDependencies || {})
  }

  const dependencyList = Object.entries(allDependencies).map(([name, ver]) => ({
    name,
    version: String(ver)
  }))

  dependencyList.sort((a, b) => a.name.localeCompare(b.name))

  return (
    <div className='p-6'>
      <div className='flex-1 flex flex-col overflow-hidden'>
        <h3 className='text-xl font-bold mb-4'>About</h3>
        <div className='flex-1 flex flex-col space-y-6 overflow-hidden'>
          <div className='space-y-2'>
            <p>
              <strong>Application Name:</strong> {appName}
            </p>
            <p>
              <strong>Version:</strong> {version}
            </p>
            <p>
              <strong>License:</strong> MIT License
            </p>
            <p className='mt-4'>For more information, visit our GitHub repository.</p>
          </div>
          <Separator className='my-4' />
          <div className='flex-1 flex flex-col min-h-0'>
            <h4 className='text-lg font-semibold mb-3'>Packages Used:</h4>
            <ScrollArea className='flex-1 border rounded-md p-3'>
              <ul className='list-disc list-inside space-y-1 text-sm'>
                {dependencyList.length > 0 ? (
                  dependencyList.map(dep => (
                    <li key={dep.name} className='flex justify-between items-center pr-2'>
                      <span className='font-medium'>{dep.name}</span>
                      <span className='text-xs ml-2'>{dep.version}</span>
                    </li>
                  ))
                ) : (
                  <p className='text-gray-500 dark:text-gray-400'>No packages found.</p>
                )}
              </ul>
            </ScrollArea>
          </div>
        </div>
      </div>
    </div>
  )
}
