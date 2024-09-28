package slices

// FilterNonNil removes every element that is nil.
//
// Parameters:
//   - slice: The slice to filter.
//
// Returns:
//   - []T: The filtered slice.
//
// NOTES: This function has side-effects, meaning that it changes the original slice.
// To avoid unintended side-effects, you may either want to use the optimized PureFilterNonNil
// or just copy the slice before applying the filter.
func FilterNonNil[T Pointer](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}

	var top int

	for i := 0; i < len(slice); i++ {
		elem := slice[i]

		if !elem.IsNil() {
			slice[top] = elem
			top++
		}
	}

	return slice[:top:top]
}
