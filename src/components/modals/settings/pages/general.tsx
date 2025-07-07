const GeneralSettings = () => {
  return (
    <div className='p-6'>
      <h3 className='text-xl font-bold mb-4'>General</h3>
      <div className='mt-6 space-y-4'>
        <div>
          <label
            htmlFor='theme-select'
            className='block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1'
          >
            App Theme
          </label>
          <select
            id='theme-select'
            className='mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white'
          >
            <option>Dark</option>
            <option>Light</option>
            <option>System</option>
          </select>
        </div>
        <div>
          <label
            htmlFor='language-select'
            className='block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1'
          >
            Language
          </label>
          <select
            id='language-select'
            className='mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md dark:bg-gray-700 dark:border-gray-600 dark:text-white'
          >
            <option>English</option>
            <option>French</option>
            <option>Spanish</option>
          </select>
        </div>
      </div>
    </div>
  )
}

export { GeneralSettings }
