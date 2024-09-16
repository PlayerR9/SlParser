package internal

// TODO: Remove this once go-commons is updated.

// SolWithLevel is a type that defines a slice of solutions with level.
type SolWithLevel[T any] struct {
	// level is the level of the solutions.
	level int

	// has_level is true if the level is set, false otherwise.
	has_level bool

	// solutions is the slice of solutions.
	solutions []T
}

// AddSolution adds a solution to the slice. Does nothing if the receiver is nil.
//
// Parameters:
//   - level: the level of the solution.
//   - solution: the solution.
//
// If the new level is lower than the current level, this function does nothing. If
// the new level is the same as the current level, the solution is added to the
// slice. Otherwise, the level is set to the new level and the slice is reset.
func (sl *SolWithLevel[T]) AddSolution(level int, solution T) {
	if sl == nil {
		return
	}

	if !sl.has_level {
		sl.level = level
		sl.has_level = true

		sl.solutions = []T{solution}
	} else {
		if level < sl.level {
			return
		}

		if level == sl.level {
			sl.solutions = append(sl.solutions, solution)
		} else {
			sl.level = level

			if len(sl.solutions) > 0 {
				for i := range sl.solutions {
					sl.solutions[i] = *new(T)
				}

				sl.solutions = sl.solutions[:0]
			}

			sl.solutions = []T{solution}
		}
	}
}

// Solutions returns the list of solutions.
//
// Returns:
//   - []T: The list of solutions.
func (sl SolWithLevel[T]) Solutions() []T {
	slice := make([]T, len(sl.solutions))
	copy(slice, sl.solutions)

	return slice
}

// Level returns the level of the solutions.
//
// Returns:
//   - int: The level of the solutions.
//   - bool: True if the level is set, false otherwise.
func (sl SolWithLevel[T]) Level() (int, bool) {
	return sl.level, sl.has_level
}

// Reset resets the level and the solutions.
func (sl *SolWithLevel[T]) Reset() {
	if sl == nil {
		return
	}

	sl.has_level = false

	if len(sl.solutions) > 0 {
		for i := range sl.solutions {
			sl.solutions[i] = *new(T)
		}

		sl.solutions = sl.solutions[:0]
	}
}
