/**
 * Generic list manager composable providing standardized add and remove operations.
 */
export function useListManager<T>(
  getList: () => T[],
  createItem?: () => T,
  beforeAdd?: () => void
) {
  function add(item?: T) {
    beforeAdd?.()
    const target = getList()
    const toPush = item !== undefined ? item : (createItem ? createItem() : undefined)
    if (toPush !== undefined) {
      target.push(toPush)
    }
  }

  function remove(index: number) {
    const target = getList()
    if (index >= 0 && index < target.length) {
      target.splice(index, 1)
    }
  }

  return {
    add,
    remove,
  }
}
