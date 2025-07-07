import { DialogTitle } from '@radix-ui/react-dialog'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'

import { Dialog, DialogContent } from '@/components/ui/dialog'
import { ScrollArea } from '@/components/ui/scroll-area'

interface Update {
  version: string
  notes: string
  pub_date: string
}

interface NewUpdateModalProps {
  update: Update | null
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

const NewUpdateModal = ({ update, open = false, onOpenChange }: NewUpdateModalProps) => {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='flex p-0 h-[250px] w-[450px] overflow-hidden'>
        <ScrollArea className='w-full h-full flex-1 p-6'>
          {update ? (
            <>
              <DialogTitle className='text-lg mb-4'>
                New Update Available: {update.version}
              </DialogTitle>
              <div className='prose dark:prose-invert max-w-none'>
                <ReactMarkdown remarkPlugins={[remarkGfm]}>
                  {update.notes || 'No release notes available for this version.'}
                </ReactMarkdown>
              </div>
            </>
          ) : (
            <>
              <DialogTitle className='text-lg'>No New Updates</DialogTitle>
              <p className='text-sm text-muted-foreground'>
                You are already using the latest version.
              </p>
            </>
          )}
        </ScrollArea>
      </DialogContent>
    </Dialog>
  )
}

export { NewUpdateModal, type NewUpdateModalProps, type Update }
