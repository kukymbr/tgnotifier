package types

import (
	"cmp"
	"fmt"
	"slices"
)

// NewNamed creates new Named instance.
func NewNamed[TNames cmp.Ordered, TValues any](
	defaultName TNames,
	nameFallbackFunc ...func(item TValues) TNames,
) *Named[TNames, TValues] {
	emptyNameFunc := func(item TValues) TNames {
		var empty TNames

		return empty
	}

	if len(nameFallbackFunc) > 0 && nameFallbackFunc[0] != nil {
		emptyNameFunc = nameFallbackFunc[0]
	}

	n := &Named[TNames, TValues]{
		items:            make(map[TNames]TValues),
		names:            make([]TNames, 0),
		nameFallbackFunc: emptyNameFunc,
	}

	n.SetDefaultName(defaultName)

	return n
}

// Named is a container for items indexed by names
// with a list of ordered names.
type Named[TNames cmp.Ordered, TValues any] struct {
	items            map[TNames]TValues
	names            []TNames
	defaultName      TNames
	nameFallbackFunc func(item TValues) TNames
}

// Add an item.
func (n *Named[TNames, TValues]) Add(name TNames, item TValues) {
	var empty TNames

	if name == empty {
		name = n.nameFallbackFunc(item)
	}

	if _, ok := n.items[name]; ok {
		return
	}

	n.items[name] = item
	n.names = append(n.names, name)

	slices.Sort(n.names)
}

// SetDefaultName sets name as a default item name.
func (n *Named[TNames, TValues]) SetDefaultName(name TNames) {
	n.defaultName = name
}

// FindByName finds item by its name.
func (n *Named[TNames, TValues]) FindByName(name TNames) (TValues, error) {
	var (
		emptyName TNames
		emptyItem TValues
	)

	if name == emptyName {
		name = n.defaultName
	}

	if name == emptyName {
		return emptyItem, fmt.Errorf("no name provided and no default name registered")
	}

	item, ok := n.items[name]
	if !ok {
		return item, fmt.Errorf("'%v' is not registered", name)
	}

	return item, nil
}

// FindByNameIndex finds item by its name index.
func (n *Named[TNames, TValues]) FindByNameIndex(index int) (TValues, error) {
	var emptyItem TValues

	if index < 0 || index >= len(n.names) {
		return emptyItem, fmt.Errorf("no name with index %d found", index)
	}

	return n.FindByName(n.names[index])
}

// GetNames returns ordered list of the items' names.
func (n *Named[TNames, TValues]) GetNames() []TNames {
	return n.names
}

// Len returns count of the items.
func (n *Named[TNames, TValues]) Len() int {
	return len(n.items)
}
