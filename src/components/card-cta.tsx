import React from 'react'

import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import { cn } from '@/lib/utils'
import { Separator } from '@/components/ui/separator'

interface CardCtaProps {
  title: string | React.ReactNode
  description?: string | React.ReactNode
  position?: 'top' | 'center' | 'bottom' | 'single'
  actionComponent?: React.ReactNode
  content?: React.ReactNode
  onClick?: () => void
  className?: string
}

const CardCta = React.forwardRef<HTMLDivElement, CardCtaProps>(
  (
    {
      title,
      description,
      position = 'single',
      actionComponent,
      content,
      onClick,
      className,
      ...props
    },
    ref
  ) => {
    const positionClasses = {
      top: 'rounded-t-lg',
      center: 'rounded-none',
      bottom: 'rounded-b-lg',
      single: ''
    }

    return (
      <Card
        ref={ref}
        className={cn(
          onClick && 'hover:bg-secondary transition-colors duration-150',
          position !== 'single' && 'rounded-none',
          positionClasses[position],
          'border-0 px-5 py-4 justify-center',
          className
        )}
        onClick={onClick}
        {...props}
      >
        <CardHeader className='p-0'>
          <CardTitle
            className={cn('font-normal text-sm', !description && 'row-span-2 self-center')}
          >
            {title}
          </CardTitle>
          {description && <CardDescription>{description}</CardDescription>}

          {actionComponent && (
            <CardAction className='w-full h-full flex justify-center items-center'>
              {actionComponent}
            </CardAction>
          )}
        </CardHeader>
        {content && <CardContent className='p-0'>{content}</CardContent>}
      </Card>
    )
  }
)

interface CardGroupProps {
  children: React.ReactNode
  merge?: boolean
  className?: string
}

const CardCtaGroup = ({ children, merge = true, className }: CardGroupProps) => {
  const validCardCtaChildren = React.Children.toArray(children).filter(
    (child): child is React.ReactElement<React.ComponentProps<typeof CardCta>> =>
      React.isValidElement(child) && child.type === CardCta
  )

  const totalChildren = validCardCtaChildren.length

  return (
    <div className={cn('grid grid-cols-1', merge ? 'gap-0' : 'gap-4', className)}>
      {validCardCtaChildren.map((child, index) => {
        const isLastChild = index === totalChildren - 1
        let position: 'top' | 'center' | 'bottom' | 'single' | undefined

        if (merge) {
          if (totalChildren === 1) {
            position = 'single'
          } else if (index === 0) {
            position = 'top'
          } else if (isLastChild) {
            position = 'bottom'
          } else {
            position = 'center'
          }
        } else {
          position = 'single'
        }

        return (
          <React.Fragment key={child.key || index}>
            {React.cloneElement(child, {
              ...child.props,
              position: position
            })}
            {merge && !isLastChild && (
              <div className='w-full bg-card px-4'>
                <Separator />
              </div>
            )}
          </React.Fragment>
        )
      })}
      {React.Children.toArray(children).filter(
        child => !(React.isValidElement(child) && child.type === CardCta)
      )}
    </div>
  )
}

export { CardCta, CardCtaGroup }
