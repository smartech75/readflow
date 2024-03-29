import React, { useCallback, useEffect } from 'react'
import { useFormState } from 'react-use-form-state'

import { Box, FormSelectField } from '../../components'
import { Theme, useLocalConfiguration } from '../../contexts/LocalConfigurationContext'

interface SwitchThemeFormFields {
  theme: Theme
}

const ThemeSwitch = () => {
  const { localConfiguration, updateLocalConfiguration } = useLocalConfiguration()

  const [formState, { select }] = useFormState<SwitchThemeFormFields>({
    theme: localConfiguration.theme,
  })

  useEffect(() => {
    const { theme } = formState.values
    if (localConfiguration.theme !== theme) {
      formState.setField('theme', localConfiguration.theme)
    }
  }, [localConfiguration, formState, updateLocalConfiguration])

  const handleThemeChange = useCallback(
    (event: React.ChangeEvent<HTMLSelectElement>) => {
      const theme = event.target.value as Theme
      updateLocalConfiguration({ ...localConfiguration, theme })
    },
    [updateLocalConfiguration, localConfiguration]
  )

  return (
    <form>
      <FormSelectField label="Theme" {...select({ name: 'theme', onChange: handleThemeChange })}>
        <option value="light">light</option>
        <option value="dark">dark</option>
        <option value="auto">auto (your system will decide)</option>
      </FormSelectField>
    </form>
  )
}

const ThemeBox = () => (
  <Box title="Theme">
    <p>Change the colors of the user interface according to your preferences.</p>
    <ThemeSwitch />
  </Box>
)

export default ThemeBox
