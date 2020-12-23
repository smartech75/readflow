import { useMutation, useQuery } from '@apollo/client'
import React, { useContext, useState } from 'react'
import { useModal } from 'react-modal-hook'
import { RouteComponentProps } from 'react-router'
import { Link } from 'react-router-dom'

import { updateCacheAfterDelete } from '../../categories/cache'
import {
  Category,
  DeleteCategoriesRequest,
  DeleteCategoriesResponse,
  GetCategoriesResponse,
} from '../../categories/models'
import { DeleteCategories, GetCategories } from '../../categories/queries'
import Button from '../../components/Button'
import ConfirmDialog from '../../components/ConfirmDialog'
import DataTable, { OnSelectedFn } from '../../components/DataTable'
import Loader from '../../components/Loader'
import Panel from '../../components/Panel'
import TimeAgo from '../../components/TimeAgo'
import { MessageContext } from '../../context/MessageContext'
import ErrorPanel from '../../error/ErrorPanel'
import { getGQLError, matchResponse } from '../../helpers'
import { usePageTitle } from '../../hooks'

const definition = [
  {
    title: 'Title',
    render: (val: Category) => (
      <Link title="Edit category" to={`/settings/categories/${val.id}`} data-test-id={val.title}>
        {val.title}
      </Link>
    ),
  },
  {
    title: 'Created',
    render: (val: Category) => <TimeAgo dateTime={val.created_at} />,
  },
  {
    title: 'Updated',
    render: (val: Category) => <TimeAgo dateTime={val.updated_at} />,
  },
]

export default ({ match }: RouteComponentProps) => {
  usePageTitle('Settings - Categories')

  const [errorMessage, setErrorMessage] = useState<string | null>(null)
  const [selection, setSelection] = useState<number[]>([])
  const { data, error, loading } = useQuery<GetCategoriesResponse>(GetCategories)
  const [deleteCategoriesMutation] = useMutation<DeleteCategoriesResponse, DeleteCategoriesRequest>(DeleteCategories)
  const { showMessage } = useContext(MessageContext)

  const onSelectedHandler: OnSelectedFn = (keys) => {
    setSelection(keys)
  }

  const deleteCategories = async (ids: number[]) => {
    try {
      const res = await deleteCategoriesMutation({
        variables: { ids },
        update: updateCacheAfterDelete(ids),
      })
      setSelection([])
      // console.log('Categories removed', res)
      if (res.data) {
        const nb = res.data.deleteCategories
        showMessage(nb > 1 ? `${nb} categories removed` : 'Category removed')
      }
    } catch (err) {
      setErrorMessage(getGQLError(err))
    }
    // eslint-disable-next-line @typescript-eslint/no-use-before-define
    hideDeleteConfirmModal()
  }

  const [showDeleteConfirmModal, hideDeleteConfirmModal] = useModal(
    () => (
      <ConfirmDialog
        title="Delete categories?"
        confirmLabel="Delete"
        onConfirm={() => deleteCategories(selection)}
        onCancel={hideDeleteConfirmModal}
      >
        Deleting a category is irreversible. Please confirm!
      </ConfirmDialog>
    ),
    [selection]
  )

  const render = matchResponse<GetCategoriesResponse>({
    Loading: () => <Loader />,
    Error: (err) => <ErrorPanel title="Unable to fetch categories">{err.message}</ErrorPanel>,
    Data: (data) => (
      <DataTable
        definition={definition}
        data={data.categories.entries}
        onSelected={onSelectedHandler}
        emptyMessage="No category yet"
      />
    ),
  })

  return (
    <Panel>
      <section>
        <header>
          <h2>Categories</h2>
          {selection.length > 0 && (
            <Button id="remove-selection" title="Remove selection" variant="danger" onClick={showDeleteConfirmModal}>
              Remove
            </Button>
          )}
          <Button
            id="add-new-category"
            title="Add new category"
            variant="primary"
            as={Link}
            to={{
              pathname: match.path + 'add',
              state: { modal: true },
            }}
          >
            Add category
          </Button>
        </header>
        <section>
          <p>Automatically organize your articles into categories using the rule engine.</p>
          {errorMessage != null && <ErrorPanel title="Unable to delete categories">{errorMessage}</ErrorPanel>}
          {render(loading, data, error)}
        </section>
      </section>
    </Panel>
  )
}
