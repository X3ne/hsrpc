import React from 'react'

interface Update {
  version: string
  notes: string
  pub_date: string
}

interface ModalContextType {
  update: Update | null
  setUpdate: React.Dispatch<React.SetStateAction<Update | null>>
  isUpdateModalOpen: boolean
  openUpdateModal: (update: Update) => void
  closeUpdateModal: () => void
}

const UpdateContext = React.createContext<ModalContextType>({
  update: null,
  setUpdate: () => {},
  isUpdateModalOpen: false,
  openUpdateModal: () => {},
  closeUpdateModal: () => {}
})

const UpdateProvider = ({ children }: { children: React.ReactNode }) => {
  const [update, setUpdate] = React.useState<Update | null>(null)
  const [isUpdateModalOpen, setIsUpdateModalOpen] = React.useState(false)

  const openUpdateModal = (update: Update) => {
    setUpdate(update)
    setIsUpdateModalOpen(true)
  }

  const closeUpdateModal = () => {
    setIsUpdateModalOpen(false)
    setUpdate(null)
  }

  React.useEffect(() => {
    if (update) {
      setIsUpdateModalOpen(true)
    }
  }, [update])

  return (
    <UpdateContext.Provider
      value={{
        update,
        setUpdate,
        isUpdateModalOpen,
        openUpdateModal,
        closeUpdateModal
      }}
    >
      {children}
    </UpdateContext.Provider>
  )
}

export { UpdateProvider, UpdateContext, type Update }
