import React, { ReactType } from 'react'
import Ink from 'react-ink'

import classes from './ButtonIcon.module.css'
import { classNames } from '../helpers'
import Icon from './Icon'
import { PropsOf } from './PropsOf'

interface ButtonIconProps {
  icon?: string
  variant?: 'default' | 'primary' | 'danger'
  loading?: boolean
  floating?: boolean
}

function ButtonIcon<Tag extends ReactType = 'button'>(props: { as?: Tag } & ButtonIconProps & PropsOf<Tag>) {
  const { as: Element = 'button', icon, variant, loading, floating, ...attrs } = props
  let className = classNames(classes.button, classes[variant], floating ? classes.floating : null)

  if (loading) {
    className = classNames(className, classes.loading)
    return (
      <button {...attrs} disabled className={className}>
        <Icon name="loop" />
      </button>
    )
  }

  return (
    <Element className={className} {...attrs} data-test={`btn-${variant}`}>
      <Icon name={icon} />
      <Ink />
    </Element>
  )
}

export default ButtonIcon
