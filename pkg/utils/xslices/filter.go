package xslices

func RemoveZero[S ~[]E, E comparable](slice S) S {
	var (
		result S
		zero   E
	)

	for _, item := range slice {
		if item != zero {
			result = append(result, item)
		}
	}

	return result
}

// RemoveZeroRef removes zero-value elements from the given slice in place.
// It modifies the original slice by filtering out elements equal to the zero value of the element type.
//
// The function:
//   - sets zero value to the removed elements to avoid memory leaks;
//   - does not allocate memory for a new slice;
//   - does not change the slice capacity;
//   - can be used with slices of any type, including pointers, structs, and interfaces.
func RemoveZeroRef[S ~[]E, E comparable](slice *S) {
	if slice == nil || *slice == nil {
		return
	}

	newSlice := (*slice)[:0]

	var zero E

	for _, item := range *slice {
		if item != zero {
			newSlice = append(newSlice, item)
		}
	}

	for i, maxI := len(newSlice), len(*slice)-1; i <= maxI; i++ {
		(*slice)[i] = zero
	}

	*slice = newSlice
}
