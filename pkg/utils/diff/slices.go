package diff

type Elements[E any] struct {
	Created []E
	Updated []E
	Deleted []E
}

// Slices compares two slices and returns the differences in terms of created,
// updated, and deleted elements.
// Parameters:
//   - newSlice: the new slice to compare.
//   - oldSlice: the old slice to compare.
//   - uniqueKey: a function to get the unique key for each element.
//   - isEqual: a function to check if two elements are equal.
//   - prepare: a function to prepare elements for comparison and update.
//     It should update the new element with the old element data.
//     For example, it can copy non-key fields from the old element to the new element
//     to ensure that the comparison and update are based on the latest data.
//
// Returns:
//   - Elements: the differences between the two slices. Order for Created, Updated
//     is preserved from the newSlice, Deleted from the oldSlice.
func Slices[S ~[]E, E any, K comparable](
	newSlice S,
	oldSlice S,
	uniqueKey func(E) K,
	isEqual func(newElt, oldElt E) bool,
	prepare func(newElt, oldElt E),
) Elements[E] {
	res := Elements[E]{
		Created: make([]E, 0, len(newSlice)),
		Updated: make([]E, 0, len(newSlice)),
		Deleted: make([]E, 0, len(oldSlice)),
	}

	oldMap := make(map[K]E, len(oldSlice))
	for _, item := range oldSlice {
		oldMap[uniqueKey(item)] = item
	}

	newMap := make(map[K]E, len(newSlice))

	for _, newItem := range newSlice {
		key := uniqueKey(newItem)
		newMap[key] = newItem

		oldItem, ok := oldMap[key]
		if !ok {
			res.Created = append(res.Created, newItem)

			continue
		}

		prepare(newItem, oldItem)

		if !isEqual(newItem, oldItem) {
			res.Updated = append(res.Updated, newItem)
		}
	}

	// Iterate over oldSlice to preserve the order of deleted elements.
	for _, oldItem := range oldSlice {
		if _, ok := newMap[uniqueKey(oldItem)]; !ok {
			res.Deleted = append(res.Deleted, oldItem)
		}
	}

	return res
}
