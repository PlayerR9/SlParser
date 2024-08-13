package generator

import (
	"cmp"
	"slices"

	itr "github.com/PlayerR9/go-commons/iterator"
)

/////////////////////////////////////////////////////

// entry is a key-value pair in an OrderedMap.
type entry[K cmp.Ordered, V any] struct {
	// Key is the key of the entry.
	Key K

	// Value is the value of the entry.
	Value V
}

// om_iterator is an iterator for an OrderedMap.
type om_iterator[K cmp.Ordered, V any] struct {
	// m is the map to iterate over.
	m *ordered_map[K, V]

	// pos is the current position in the iterator.
	pos int
}

// Apply implements the iterator.Iterable interface.
//
// The variable passed to the iterator function is of type entry[K, V].
func (i *om_iterator[K, V]) Apply(fn itr.IteratorFunc) error {
	// dbg.AssertNil(i.m, "i.m")

	if i.pos >= len(i.m.keys) {
		return itr.ErrExausted
	}

	key := i.m.keys[i.pos]
	i.pos++

	val := i.m.values[key]
	// dbg.AssertOk(ok, "i.m.values[%s]", strconv.Quote(lustr.GoStringOf(key)))

	err := fn(entry[K, V]{
		Key:   key,
		Value: val,
	})
	if err != nil {
		return err
	}

	return nil
}

// Restart implements the common.Iterater interface.
func (i *om_iterator[K, V]) Reset() {
	i.pos = 0
}

/* // new_om_iterator creates a new OMIterator.
//
// Parameters:
//   - m: The map to iterate over.
//
// Returns:
//   - *OMIterator: A pointer to the newly created OMIterator. Nil if m is nil.
func new_om_iterator[K cmp.Ordered, V any](m *ordered_map[K, V]) *om_iterator[K, V] {
	if m == nil {
		return nil
	}

	return &om_iterator[K, V]{
		m:   m,
		pos: 0,
	}
} */

// ordered_map is a map that is ordered by the keys.
type ordered_map[K cmp.Ordered, V any] struct {
	// values is a map of the values in the map.
	values map[K]V

	// keys is a slice of the keys in the map.
	keys []K
}

// Iterator implements the iterator.Iterater interface.
func (m ordered_map[K, V]) Iterator() itr.Iterable {
	return &om_iterator[K, V]{
		m:   &m,
		pos: 0,
	}
}

// new_ordered_map creates a new OrderedMap.
//
// Returns:
//   - *OrderedMap: A pointer to the newly created OrderedMap.
//     Never returns nil.
func new_ordered_map[K cmp.Ordered, V any]() *ordered_map[K, V] {
	return &ordered_map[K, V]{
		values: make(map[K]V),
		keys:   make([]K, 0),
	}
}

// add adds a key-value pair to the map.
//
// Parameters:
//   - key: The key to add.
//   - value: The value to add.
//   - force: If true, the value will be added even if the key already exists. If
//     false, the value will not be added if the key already exists.
//
// Returns:
//   - bool: True if the value was added to the map, false otherwise.
func (m *ordered_map[K, V]) add(key K, value V, force bool) bool {
	pos, ok := slices.BinarySearch(m.keys, key)

	if !ok {
		m.keys = slices.Insert(m.keys, pos, key)
	}

	if ok && !force {
		return false
	}

	m.values[key] = value

	return true
}

// size is a method that returns the number of keys in the map.
//
// Returns:
//   - int: The number of keys in the map.
func (m ordered_map[K, V]) size() int {
	return len(m.keys)
}

// Map is a method that returns the map of the values in the map.
//
// Returns:
//   - map[K]V: The map of the values in the map. Never returns nil.
func (m ordered_map[K, V]) Map() map[K]V {
	return m.values
}

// Keys is a method that returns the keys in the map.
//
// Returns:
//   - []K: The keys in the map.
func (m ordered_map[K, V]) Keys() []K {
	return m.keys
}
