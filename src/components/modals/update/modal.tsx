import { useEffect, useState } from 'react'
import { DialogTitle } from '@radix-ui/react-dialog'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { error } from '@tauri-apps/plugin-log'
import { Channel, invoke } from '@tauri-apps/api/core'
import { toast } from 'sonner'

import { Dialog, DialogContent } from '@/components/ui/dialog'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { useUpdate } from '@/hooks/use-update'

type DownloadEvent =
  | {
      event: 'started'
    }
  | {
      event: 'progress'
      data: {
        downloaded: number
        contentLength: number | null
      }
    }
  | {
      event: 'finished'
    }

const NewUpdateModal = () => {
  const { update, isUpdateModalOpen, closeUpdateModal } = useUpdate()

  const onEvent = new Channel<DownloadEvent>()
  const [progress, setProgress] = useState<{
    state?: 'started' | 'progress' | 'finished'
    downloaded: number
    content_length: number | null
    percent: number
  }>({ downloaded: 2e6, content_length: 1e7, percent: 0 })

  useEffect(() => {
    if (isUpdateModalOpen) {
      setProgress({ downloaded: 2e6, content_length: 1e7, percent: 0 })
    }
  }, [isUpdateModalOpen, update])

  const bytesToMB = (bytes: number) => {
    return (bytes / 1024 / 1024).toFixed(2)
  }

  onEvent.onmessage = message => {
    console.debug('Download event received:', message)
    if (message.event === 'started') {
      setProgress({ state: 'started', downloaded: 0, content_length: null, percent: 0 })
    } else if (message.event === 'progress') {
      setProgress({
        state: 'progress',
        downloaded: message.data.downloaded,
        content_length: message.data.contentLength,
        percent: message.data.contentLength
          ? (message.data.downloaded / message.data.contentLength) * 100
          : 0
      })
    } else if (message.event === 'finished') {
      setProgress({
        state: 'finished',
        downloaded: 0,
        content_length: 0,
        percent: 100
      })
      toast.success('Update downloaded successfully!')
    }
  }

  return (
    <Dialog open={isUpdateModalOpen} onOpenChange={closeUpdateModal}>
      <DialogContent className='flex flex-col h-[280px] w-[450px] overflow-hidden'>
        {update ? (
          <>
            <DialogTitle className='text-lg'>New Update Available: {update.version}</DialogTitle>
            <div className='prose dark:prose-invert max-w-none'>
              <ScrollArea className='w-full h-full flex-1'>
                <ReactMarkdown remarkPlugins={[remarkGfm]}>
                  {update.notes || 'No release notes available for this version.'}
                </ReactMarkdown>
              </ScrollArea>
            </div>
            {progress.state ? (
              <div className='w-full mt-4'>
                <Progress value={progress.percent} className='w-full' />
                <p className='text-sm text-muted-foreground mt-2'>
                  Downloading update... {bytesToMB(progress.downloaded)} MB
                  {progress.content_length ? ` / ${bytesToMB(progress.content_length)} MB` : ''}
                </p>
              </div>
            ) : (
              <Button
                className='self-end justify-self-end mt-auto'
                onClick={() => {
                  invoke('download_and_install_update', {
                    onDownload: onEvent
                  }).catch(e => {
                    error('Failed to install update:', {
                      keyValues: {
                        error: e instanceof Error ? e.message : String(e)
                      }
                    })
                    toast.error(
                      'Failed to install update. Please check logs and report this issue.'
                    )
                  })
                }}
              >
                Install Update
              </Button>
            )}
          </>
        ) : (
          <>
            <DialogTitle className='text-lg'>No New Updates</DialogTitle>
            <p className='text-sm text-muted-foreground'>
              You are already using the latest version.
            </p>
          </>
        )}
      </DialogContent>
    </Dialog>
  )
}

export { NewUpdateModal }
