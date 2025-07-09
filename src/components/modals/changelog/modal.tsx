import { useState, useEffect, useMemo } from 'react'
import { DialogTitle } from '@radix-ui/react-dialog'
import { Clock } from 'iconsax-reactjs'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { useQuery } from '@tanstack/react-query'

import { Dialog, DialogContent } from '@/components/ui/dialog'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarGroup,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarInset,
  SidebarProvider
} from '@/components/ui/sidebar'

interface ChangelogEntry {
  version: string
  content: string
  markdown: string
}

interface ChangelogModalProps {
  changelogUrl: string
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

const parseChangelog = (markdown: string): ChangelogEntry[] => {
  const entries: ChangelogEntry[] = []
  const versionHeaderRegex = /^##\s*\[(\d+\.\d+\.\d+(?:-\S+)?(?:\]\(.*?\))?)\s*\][^\n]*(?:\n|$)/gm

  const matches = [...markdown.matchAll(versionHeaderRegex)]

  for (let i = 0; i < matches.length; i++) {
    const match = matches[i]
    const fullHeaderLine = match[0]
    const version = match[1].trim()
    const startIndex = match.index!

    let endIndex
    if (i + 1 < matches.length) {
      endIndex = matches[i + 1].index!
    } else {
      endIndex = markdown.length
    }

    const fullMarkdown = markdown.substring(startIndex, endIndex).trim()

    const contentStartIndex = startIndex + fullHeaderLine.length
    const content = markdown.substring(contentStartIndex, endIndex).trim()

    entries.push({
      version: version,
      content: content,
      markdown: fullMarkdown
    })

    if (entries.length >= 6) {
      break
    }
  }

  return entries
}

const ChangelogModal = ({ changelogUrl, open = false, onOpenChange }: ChangelogModalProps) => {
  const {
    data: rawChangelogContent,
    isLoading,
    isError,
    error
  } = useQuery<string, Error>({
    queryKey: ['changelog', changelogUrl],
    queryFn: async () => {
      const response = await fetch(changelogUrl)
      if (!response.ok) {
        throw new Error(`Failed to fetch changelog: ${response.statusText}`)
      }
      const text = await response.text()
      return text
    },
    staleTime: 1000 * 60 * 60,
    refetchOnWindowFocus: false
  })

  const changelogEntries = useMemo(() => {
    if (rawChangelogContent) {
      return parseChangelog(rawChangelogContent)
    }
    return []
  }, [rawChangelogContent])

  const [activeVersion, setActiveVersion] = useState<string | null>(null)

  useEffect(() => {
    if (changelogEntries.length > 0 && activeVersion === null) {
      const firstVersion = changelogEntries[0].version
      setActiveVersion(firstVersion)
    }
  }, [changelogEntries, activeVersion])

  const currentVersionContent = useMemo(() => {
    const entry = changelogEntries.find(entry => entry.version === activeVersion)
    return entry ? entry.markdown : 'No changelog entry found for this version.'
  }, [activeVersion, changelogEntries])

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='flex p-0 h-[400px] w-[800px] overflow-hidden'>
        <DialogTitle className='sr-only'>Changelog</DialogTitle>{' '}
        <SidebarProvider className='h-[400px] min-h-0'>
          <Sidebar className='min-w-[200px]'>
            <SidebarHeader className='px-4 mt-4'>
              <h2 className='text-xl'>Changelog</h2>
            </SidebarHeader>
            <SidebarContent>
              {isLoading && <p className='p-4 text-sm'>Loading versions...</p>}
              {isError && <p className='p-4 text-sm'>Error loading versions: {error?.message}</p>}
              {!isLoading && !isError && changelogEntries.length === 0 && (
                <p className='p-4 text-sm'>No changelog entries found.</p>
              )}
              {!isLoading && !isError && changelogEntries.length > 0 && (
                <SidebarGroup>
                  <SidebarMenu>
                    {changelogEntries.map(entry => (
                      <SidebarMenuItem key={entry.version}>
                        <SidebarMenuButton
                          tooltip={`View changes for ${entry.version}`}
                          isActive={entry.version === activeVersion}
                          onClick={() => setActiveVersion(entry.version)}
                        >
                          <Clock size={18} />
                          <span>{entry.version}</span>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    ))}
                  </SidebarMenu>
                </SidebarGroup>
              )}
            </SidebarContent>
          </Sidebar>

          <ScrollArea className='w-full h-full flex-1 p-6'>
            <SidebarInset>
              {isLoading && <p className='text-center'>Loading changelog content...</p>}
              {isError && <p className='text-center'>Error: {error?.message}</p>}
              {!isLoading && !isError && (
                <div className='prose dark:prose-invert max-w-none'>
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>{currentVersionContent}</ReactMarkdown>
                </div>
              )}
            </SidebarInset>
          </ScrollArea>
        </SidebarProvider>
      </DialogContent>
    </Dialog>
  )
}

export { ChangelogModal }
