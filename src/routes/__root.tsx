import { Layout } from '@/components/layout/layout'
import { createRootRoute, Outlet } from '@tanstack/react-router'

export const Route = createRootRoute({
  component: () => (
    <>
      <Layout>
        <Outlet />
      </Layout>
    </>
  )
})
