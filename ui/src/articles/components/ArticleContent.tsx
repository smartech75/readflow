import React, { useEffect, useRef } from 'react'

import { Article } from '../models'
import styles from './ArticleContent.module.css'

interface Props {
  article: Article
}

export default ({ article }: Props) => {
  const contentRef = useRef<any>(null)

  var cssLink = document.createElement('link')
  cssLink.href = process.env.PUBLIC_URL + '/readable.css'
  cssLink.rel = 'stylesheet'
  cssLink.type = 'text/css'
  var script = document.createElement('script')
  script.setAttribute('type', 'text/javascript')
  script.setAttribute('src', process.env.PUBLIC_URL + '/readable.js')

  useEffect(() => {
    let ifrm = contentRef.current
    ifrm = ifrm.contentWindow || ifrm.contentDocument.document || ifrm.contentDocument
    ifrm.document.open()
    ifrm.document.write(article.html)
    ifrm.document.head.appendChild(cssLink)
    ifrm.document.head.appendChild(script)
    ifrm.document.close()
  }, [article])

  return (
    <article className={styles.content}>
      <iframe ref={contentRef} />
    </article>
  )
}
