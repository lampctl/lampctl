/**
 * Sort an array of objects by property
 * @param {*} a array of objects
 * @param {*} prop name of property in the objects
 * @returns sorted array
 */
function sortByProp(a, prop) {
  return a.sort((a, b) => {
    if (a[prop] < b[prop]) { return -1 }
    if (b[prop] < a[prop]) { return 1 }
    return 0
  })
}

export {
  sortByProp
}
